package garage

import (
	"github.com/nsf/termbox-go"
)

type Gui struct {
}

func (g* Gui) Start() {
	termbox.Init()
}

func (g* Gui) Stop() {
	termbox.Close()
}

func (g* Gui) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (g* Gui) PollEvent() termbox.Event {
	return termbox.PollEvent()
}

func (g* Gui) PrintString(x, y int, s string) {
	for _, r := range s {
		fontColor := termbox.ColorDefault
		bgColor := termbox.ColorDefault
		termbox.SetCell(x, y, r, fontColor, bgColor)
		x ++
	}
}

func (g* Gui) Flush() {
	termbox.Flush()
}
