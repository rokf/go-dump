package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/gorilla/websocket"
)

func main() {
	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial("ws://echo.websocket.org", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	s, _ := tcell.NewScreen()
	s.Init()

	s.Clear()

	quit := make(chan struct{})
	received := make(chan []byte)

	msg := []rune{}
	allMessages := []string{}

	go func() {
		for {
			_, m, e := conn.ReadMessage()
			if e != nil {
				break
			}
			received <- m
		}
	}()

	events := make(chan tcell.Event, 100)
	var event tcell.Event

	go func() {
		for {
			events <- s.PollEvent()
		}
	}()

	go func() {
		for {
			select {
			case event = <-events:
				switch ev := event.(type) {
				case *tcell.EventKey:
					switch ev.Key() {
					case tcell.KeyEsc:
						close(quit)
						return
					case tcell.KeyCtrlU:
						msg = msg[:0]
					case tcell.KeyEnter:
						conn.WriteMessage(websocket.TextMessage, []byte(string(msg)))
						msg = msg[:0]
					case tcell.KeyRune:
						msg = append(msg, ev.Rune())
					}
				case *tcell.EventResize:
					s.Sync()
				}
			case <-time.After(time.Millisecond * 50):
			}

			// draw
			s.Clear()
			_, h := s.Size()
			for mi, m := range allMessages {
				for ri, r := range m {
					s.SetContent(ri+1, mi+1, r, nil, tcell.StyleDefault)
				}
			}
			for ri, r := range msg {
				s.SetContent(ri+1, h-1, r, nil, tcell.StyleDefault)
			}
		}

	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case r := <-received:
			allMessages = append(allMessages, string(r))
			_, h := s.Size()
			if len(allMessages) > (h - 3) {
				allMessages = allMessages[1:]
			}
		case <-time.After(time.Millisecond * 50):
		}
		s.Show()
	}

	s.Fini()
}
