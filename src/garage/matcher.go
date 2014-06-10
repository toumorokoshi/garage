package garage

import (
	"os"
	"bufio"
	"fmt"
	"github.com/toumorokoshi/go-fuzzy/fuzzy"
  "github.com/nsf/termbox-go"
	"log"
	"strings"
)

const (
	AUTOCOMPLETE_PREFIX = "gfind:"
	MESSAGE_ARGUMENT_SEPARATOR = ";"
)

// represents a single entry in the autocomplete library
type CompleteEntry struct {
	Message string
	Command string
}


type GarageMatcher struct {
	candidates []CompleteEntry
	matcher fuzzy.Matcher
	completed bool
	gui *Gui
}

func draw(g *Gui, input string, candidates fuzzy.Matches) {
	// draw the garagecomplete
	g.Clear()
	g.PrintString(0, 0, "Garage Complete")
	completionsFound := fmt.Sprintf(
		"%d total completions found",
		len(candidates),
	)
	g.PrintString(0, 1, completionsFound)
	g.PrintString(0, 2, "I want to: " + input)

	row := 3
	for i := range candidates {
		entryString := fmt.Sprintf("%d: %s ; %s",
			i + 1,
			candidates[i].Value,
			candidates[i].Data["Command"],
		)
		g.PrintString(row, 0, entryString)
		row++
	}
}

func (gm *GarageMatcher) Start() {
	gm.completed = false

	window := &Gui{}
	gm.gui = window

	input := ""
	command := ""
	for !gm.completed {
		currentCandidates := gm.matcher.ClosestList(input, 10)
		draw(window, input, currentCandidates)

		switch event := window.PollEvent(); event.Type {
		case termbox.EventKey:

			switch event.Key {

			// backspace
			case termbox.KeyBackspace:
				if len(input) > 0 {
					input = input[0: len(input) - 1]
				}
				continue

			// shutdown commands
			case termbox.KeyEsc:
				gm.Stop()
				continue

			case termbox.KeyCtrlC:
				gm.Stop()
				continue

			// succesful command
			case termbox.KeyEnter:
				gm.completed = true
				command = currentCandidates[0].Data["Command"]
				continue

			default:
				input += string(event.Key)
			}
		}
	}
	gm.Stop()
	fmt.Println(command)
	writeToFileDescriptor(command)
}

func writeToFileDescriptor(message string) {
	// write to file descriptor three, a hack to be able
	// to execute functions and also have an ncurses interface
	file := os.NewFile(3, "mythicalthree")
	file.Write([]byte(message))
}

func (gm* GarageMatcher) Stop() {
	gm.gui.Stop()
	gm.completed = true
}

func NewGarageMatcher(completeEntries []CompleteEntry) *GarageMatcher {
	return &GarageMatcher{
		completeEntries,
		createMatcher(completeEntries),
		false,
		nil,
	}
}

func createMatcher(candidateEntries []CompleteEntry) fuzzy.Matcher {
	candidates := make([]fuzzy.MatchStruct, len(candidateEntries), len(candidateEntries))
	for i := range candidateEntries {
		candidates[i] = createMatchStruct(candidateEntries[i])
	}
	return fuzzy.NewMatcher(candidates)
}

func createMatchStruct(completeEntry CompleteEntry) fuzzy.MatchStruct {
	return fuzzy.MatchStruct {
		completeEntry.Message,
		map[string]string {
			"Command": completeEntry.Command,
		},
	}
}

func CreateMatcherFromRepository(repository *GarageRepository) *GarageMatcher {
	// create a matcher from a garage repository
	completeEntries := make(
		[]CompleteEntry,
		0,
		len(repository.Scripts),
	)
	for _, script := range repository.Scripts {
		fullScriptPath := strings.Join([]string{repository.RootPath, script}, "/")
		completeEntries = addCompletionFromScript(completeEntries, fullScriptPath)
	}
	return NewGarageMatcher(completeEntries)
}

func addCompletionFromScript(completeEntries []CompleteEntry, script string) []CompleteEntry {
	// add the completions found in the script
	file, err := os.Open(script)
	if err != nil {
		log.Fatal("addCompletionFromScript: ", err)
	}
	reader := bufio.NewReader(file)
	isPrefix := true
	line := ""

	for isPrefix && err == nil {
		line, err = reader.ReadString('\n')
		completeEntry := getCompleteEntryFromString(line)
		if (completeEntry != nil) {
			completeEntries = append(completeEntries, *completeEntry)
		}
	}
	return completeEntries
}

func getCompleteEntryFromString(line string) *CompleteEntry {
	splitString := strings.Split(line, AUTOCOMPLETE_PREFIX)
	if len(splitString) == 1 {
		return nil
	}
	entryString := strings.Split(splitString[1], MESSAGE_ARGUMENT_SEPARATOR)
	if len(entryString) == 1 {
		return &CompleteEntry{strings.TrimSpace(entryString[0]), ""}
	} else {
		return &CompleteEntry{
			strings.TrimSpace(entryString[0]),
			strings.TrimSpace(entryString[1]),
		}
	}
}
