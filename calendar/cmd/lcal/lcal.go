package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var event tcell.Event
var quit chan struct{}
var events chan tcell.Event
var msg []rune
var db *gorm.DB
var data []entry
var commandParser *regexp.Regexp
var selectedEntry int

type entry struct {
	ID   int
	Date time.Time
	Text string
}

func main() {
	db, _ = gorm.Open("sqlite3", "/tmp/lcal.db")
	defer db.Close()

	if !db.HasTable("entries") {
		db.Table("entries").CreateTable(&entry{})
	}

	commandParser = regexp.MustCompile(`(\d+) (\d+) (\d+) (.+)`)
	quit = make(chan struct{})
	events = make(chan tcell.Event, 100)
	msg = []rune{}
	selectedEntry = 0
	db.Order("date").Find(&data)

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Clear()

	go func() {
		for {
			events <- s.PollEvent()
		}
	}()

loop:
	for {
		select {
		case event = <-events:
			switch ev := event.(type) {
			case *tcell.EventKey:
				handleKeyEvent(ev)
			case *tcell.EventResize:
				s.Sync()
			}
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
		draw(s)
		s.Show()
	}

	s.Fini()
}

func draw(s tcell.Screen) {
	s.Clear()
	drawCurrentDate(s)
	drawData(s)
	drawCommand(s)
	drawSelectedEntry(s)
}

func makeDateString(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d %02d %02d", year, int(month), day)
}

func drawSelectedEntry(s tcell.Screen) {
	if len(data) == 0 {
		return
	}
	w, _ := s.Size()
	s.SetContent(w-1, selectedEntry+1, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorDefault))
}

func drawData(s tcell.Screen) {
	for ei, e := range data {
		newLine := fmt.Sprintf("%v %v", makeDateString(e.Date), e.Text)
		drawString(s, []rune(newLine), 0, ei+1, tcell.ColorDefault, tcell.ColorDefault)
	}
}

func drawString(s tcell.Screen, str []rune, x, y int, bgcolor tcell.Color, fgcolor tcell.Color) {
	for ri, r := range str {
		s.SetContent(x+ri, y, r, nil, tcell.StyleDefault.Background(bgcolor).Foreground(fgcolor))
	}
}

func drawCurrentDate(s tcell.Screen) {
	w, _ := s.Size()
	dateStr := makeDateString(time.Now()) + " today"
	drawString(s, []rune(dateStr+strings.Repeat(" ", w-len(dateStr))), 0, 0, tcell.ColorWhite, tcell.ColorBlack)
}

func drawCommand(s tcell.Screen) {
	_, h := s.Size()

	drawString(s, []rune(strings.Replace(string(msg), " ", "·", -1)),
		0, h-1, tcell.ColorDefault, tcell.ColorDefault)
}

func handleKeyEvent(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		close(quit)
	case tcell.KeyRune:
		msg = append(msg, ev.Rune())
	case tcell.KeyBackspace:
	case tcell.KeyBackspace2:
		if len(msg) != 0 {
			msg = msg[:len(msg)-1]
		}
	case tcell.KeyCtrlU:
		msg = msg[:0]
	case tcell.KeyDown:
		if selectedEntry < (len(data) - 1) {
			selectedEntry = selectedEntry + 1
		}
	case tcell.KeyUp:
		if selectedEntry > 0 {
			selectedEntry = selectedEntry - 1
		}
	case tcell.KeyCtrlD:
		if len(data) == 0 {
			break
		}
		db.Where("id = ?", data[selectedEntry].ID).Delete(entry{})
		db.Order("date").Find(&data)
		selectedEntry = 0
	case tcell.KeyEnter:
		if len(msg) == 0 {
			break
		}

		messageData := commandParser.FindStringSubmatch(string(msg))
		if len(messageData) != 5 {
			break
		}
		year, err := strconv.Atoi(messageData[1])
		month, err := strconv.Atoi(messageData[2])
		day, err := strconv.Atoi(messageData[3])
		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		if err == nil {
			newEntry := entry{
				Text: string(messageData[4]),
				Date: t,
			}

			if result := db.Create(&newEntry); result.Error != nil {
				panic(result.Error)
			}

			db.Order("date").Find(&data)
		}

		msg = msg[:0]
	case tcell.KeyCtrlE:
		if len(msg) != 0 || len(data) == 0 {
			break
		}
		msg = []rune(fmt.Sprintf("%v %v", makeDateString(data[selectedEntry].Date), data[selectedEntry].Text))
		db.Where("id = ?", data[selectedEntry].ID).Delete(entry{})
		db.Order("date").Find(&data)
		selectedEntry = 0
	}
}
