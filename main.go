package main

import "log"
import "os"
import "github.com/docopt/docopt-go"
import "garage"
import "os/signal"

const (
	DEFAULT_DIRECTORY = "."
)

func main() {
	usage := `
garage

Usage:
  garage [-d <directory>]

Options:
  -d <directory>, --directory <directory>    The directory to retrieve scripts from.
                                             if not passed in, the current directory
                                             will be used.
  `

	args, _ := docopt.Parse(usage, nil, true, "garage 0.1", false)

	directory := DEFAULT_DIRECTORY
	if str, ok := args["--directory"].(string); ok {
		directory = str
	}

	repository := garage.LoadGarageRepository(directory)
	garageMatcher := garage.CreateMatcherFromRepository(repository)

	gui := garage.NewDefaultGui()
	gui.Print("hello world!")
	gui.Clear(5)
	gui.GetChar()

	// start the listener
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(garageMatcher *garage.GarageMatcher){
		for sig := range(c) {
			log.Print(sig)
			garageMatcher.Stop()
			os.Exit(0)
		}
	}(garageMatcher)

	garageMatcher.Start()

}
