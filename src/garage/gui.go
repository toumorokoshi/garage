package garage

import (
	"os"
)

type Gui struct {
	input *os.File
	output *os.File
}

func NewDefaultGui() *Gui {
	return &Gui{os.Stdin, os.Stdout}
}

func (g* Gui) Print(message string) {
	g.output.WriteString(message)
}

func (g* Gui) Clear(count int) {
	for i := 0; i < count; i++ {
		g.output.WriteString("\b")
	}

	for i := 0; i < count; i++ {
		g.output.WriteString(" ")
	}
}

func (g* Gui) GetChar() byte {
	readChar := make([]byte, 1)
	g.input.Read(readChar)
	return readChar[0]
}
