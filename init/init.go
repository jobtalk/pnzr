package init

import (
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

var LANG = "JP"

var sections = map[string]map[string]Section{
	"JP": JPSections,
}

type Section interface {
	Answer() string
	Message() string
	ToggleCursor()
}

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)
	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}

func drawString(s string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawLine(0, 0, "Press ESC to exit.")
	drawLine(0, 1, s)
	termbox.Flush()
}

func draw(d Section) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawLine(0, 0, "Press ESC to exit.")

	for i, l := range strings.Split(d.Message(), "\n") {
		drawLine(0, i+1, l)
	}

	termbox.Flush()
}

func RunInit() {
	box := NewPolarQuestionBox("test")
	var fin = make(chan bool)
	//点滅を制御する
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-t.C:
				box.ToggleCursor()
				draw(box)
			case <-fin:
			default:
			}
		}
		t.Stop()
	}()
	//box := sections["JP"]["generateQuestinType"].(*SelectBox)
	draw(box)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				fin <- true
				return
			case termbox.KeyCtrlC:
				fin <- true
				return
			case termbox.KeyEnter:
				drawString(box.Answer())
			case termbox.KeyBackspace:
				box.BS()
				draw(box)
			case termbox.KeyBackspace2:
				box.BS()
				draw(box)
			default:
				box.Add(ev.Ch)
				draw(box)
			}
		default:
			draw(box)
		}
	}
}
