package main

import "fmt"
import "log"
import "github.com/docopt/docopt-go"
import "code.google.com/p/goncurses"

func main() {
	usage := `
garage

Usage:
  garage
  `

	arguments, _ := docopt.Parse(usage, nil, true, "garage 0.1", false)
	fmt.Println(arguments)
	// ncurses
	screen, err := goncurses.Init()
	if err != nil {
		log.Fatal("init: ", err)
	}
	defer goncurses.End()

	screen.Print("Hello, World")
	screen.MovePrint(3, 0, "Press any key to continue")
	screen.Refresh()
	screen.GetChar()
}
