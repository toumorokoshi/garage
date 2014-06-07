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
	AUTOCOMPLETEPREFIX = "gfind:"
	MESSAGEARGUMENTSSEPARATOR = ";"
)

// represents a single entry in the autocomplete library
type CompleteEntry struct {
	Message string
	Command string
}


type GarageMatcher struct {
	candidates []CompleteEntry
	matcher fuzzy.Matcher
}

func (gm *GarageMatcher) Start() {
	screen, err := goncurses.Init()
	if err != nil {
		log.Fatal("init: ", err)
	}
	defer goncurses.End()

	// we handle printing characters ourselves
	goncurses.Echo(false)

	screen.Print("Garage Complete")
	i := 1

	screen.Move(i, 0)
	screen.Printf("%d completions found", len(gm.candidates))
	i++

	screen.MovePrint(i, 0, "I want to: ")
	input := ""
	completed := false
	for(!completed) {
		char := screen.GetChar()
		switch char {
		case goncurses.KEY_RETURN:
			completed = true
			continue
		case 27:
			// escape key
			completed = true
			continue
		}
		input += string(char)
		screen.MovePrint(i, 12, input)
		currentCandidates := gm.matcher.ClosestList(input, 10)
		printMatches(screen, currentCandidates, i + 1)
		screen.Move(i, 12 + len(input))
	}
}

func NewGarageMatcher(completeEntries []CompleteEntry) *GarageMatcher {
	return &GarageMatcher{
		completeEntries,
		createMatcher(completeEntries),
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

func getAutoCompleteCandidates() []CompleteEntry{
	return []CompleteEntry{
		CompleteEntry { "find all files in a directory", "find -r *" },
		CompleteEntry { "write to a file", "cat"},
		CompleteEntry { "copy a directory", "cp -r"},
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
	splitString := strings.Split(line, AUTOCOMPLETEPREFIX)
	if len(splitString) == 1 {
		return nil
	}
	entryString := strings.Split(splitString[1], MESSAGEARGUMENTSSEPARATOR)
	if len(entryString) == 1 {
		return &CompleteEntry{strings.TrimSpace(entryString[0]), ""}
	} else {
		return &CompleteEntry{
			strings.TrimSpace(entryString[0]),
			strings.TrimSpace(entryString[1]),
		}
	}
}
