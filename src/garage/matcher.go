package garage

import (
	"os"
	"bufio"
	"fmt"
	"github.com/toumorokoshi/go-fuzzy/fuzzy"
	"code.google.com/p/goncurses"
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
	screen *goncurses.Screen
}

func (gm* GarageMatcher) initWindow() *goncurses.Window {

	/* input, err := os.OpenFile("/dev/tty", 0, 0444)
	if err != nil {
		log.Fatal("initWindow: ", err)
	}
	output, err := os.OpenFile("/dev/tty", 0, 0222)
	if err != nil {
		log.Fatal("initWindow: ", err)
	}

	screen, err := goncurses.NewTerm("", output, input)
	if err != nil {
		log.Fatal("initWindow: ", err)
	}

	screen.Set()
	gm.screen = screen
	return goncurses.StdScr()
  */

	window, err := goncurses.Init()
	if err != nil {
		log.Fatal("initWindow:", err)
	}
	return window
}

func (gm *GarageMatcher) Start() {
	gm.completed = false

	window := gm.initWindow()

	// we handle printing characters ourselves
	goncurses.Echo(false)

	window.Print("Garage Complete")
	i := 1

	window.Move(i, 0)
	window.Printf("%d total completions found", len(gm.candidates))
	i++

	input := ""
	command := ""
	for(!gm.completed) {
		currentCandidates := gm.matcher.ClosestList(input, 10)
		printMatches(window, currentCandidates, i + 1)
		window.Move(i, 0)
		window.ClearToEOL()
		window.MovePrint(i, 0, "I want to: " + input)

		char := window.GetChar()
		switch char {
		case 127:
			// backspace
			if len(input) > 0 {
				input = input[0: len(input) - 1]
			}
			continue
		case goncurses.KEY_RETURN:
			gm.completed = true
			command = currentCandidates[0].Data["Command"]
			continue
		case 27:
			// escape key
			gm.Stop()
			continue
		}
		input += string(char)
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
	/* if (gm.screen != nil) {
		gm.screen.End()
		gm.screen.Delete()
	} */
	goncurses.End()
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

func printMatches(screen *goncurses.Window, candidates fuzzy.Matches, rowToDraw int) {
	for i := range candidates {
		screen.Move(rowToDraw, 0)
		screen.ClearToEOL()
		entryString := fmt.Sprintf("%d: %s ; %s",
			i + 1,
			candidates[i].Value,
			candidates[i].Data["Command"],
		)
		screen.MovePrint(rowToDraw, 0, entryString)
		rowToDraw++
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
