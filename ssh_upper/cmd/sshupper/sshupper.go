package main

import (
	"bufio"
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/tidwall/transform"
)

func toUpper(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	return transform.NewTransformer(func() ([]byte, error) {
		c, _, err := br.ReadRune()
		if err != nil {
			return nil, err
		}
		return []byte(strings.ToUpper(string([]rune{c, '\r'}))), nil
	})
}

func main() {
	ssh.Handle(func(s ssh.Session) {
		io.Copy(s, toUpper(s))
	})
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
