package main

import "github.com/docopt/docopt-go"
import "garage"

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
	garageMatcher.Start()
}
