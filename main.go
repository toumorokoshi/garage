package main

import "fmt"
import "github.com/docopt/docopt-go"
import "garage"

func main() {
	usage := `
garage

Usage:
  garage
  `

	arguments, _ := docopt.Parse(usage, nil, true, "garage 0.1", false)
	fmt.Println(arguments)

	repository := garage.LoadGarageRepository("/home/tsutsumi/workspace/sub/libexec")
	garageMatcher := garage.CreateMatcherFromRepository(repository)
	garageMatcher.Start()
}
