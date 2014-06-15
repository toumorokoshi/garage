package main

import "github.com/docopt/docopt-go"
import "garage"

const (
	DEFAULT_DIRECTORY = "."
	DEFAULT_NAME = "Garage"
)

func main() {
	usage := `
garage

Usage:
  garage [-d <directory> -n <name>]

Options:
  -d <directory>, --directory <directory>    The directory to retrieve scripts from.
                                             if not passed in, the current directory
                                             will be used.
  -n <name>, --name <name>                   The name of the garage repository.
                                             "Garage" is the default
  `

	args, _ := docopt.Parse(usage, nil, true, "garage 0.1", false)

	directory := DEFAULT_DIRECTORY
	name := DEFAULT_NAME

	if str, ok := args["--directory"].(string); ok {
		directory = str
	}

	if str, ok := args["--name"].(string); ok {
		name = str
	}

	repository := garage.LoadGarageRepository(directory)
	garageMatcher := garage.CreateMatcherFromRepository(repository, name)
	garageMatcher.Start()
}
