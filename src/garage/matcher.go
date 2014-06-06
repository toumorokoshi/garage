package garage

import "fmt"
import "github.com/toumorokoshi/go-fuzzy/fuzzy"
import "code.google.com/p/goncurses"
import "log"

// represents a single entry in the autocomplete library
type CompleteEntry struct {
	Message string
	CommandArguments string
}


type GarageMatcher struct {
	screen* goncurses.Window
	candidates []string
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

	candidates := getAutoCompleteCandidates()
	gm.matcher = createMatcher(candidates)
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
			"CommandArguments": completeEntry.CommandArguments,
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
			candidates[i].Data["CommandArguments"],
		)
		screen.MovePrint(rowToDraw, 0, entryString)
		rowToDraw++
	}
}

func matchCandidates(candidates []string, input string, limit int) []string {
	return candidates
}

func getAutoCompleteCandidates() []CompleteEntry{
	return []CompleteEntry{
		CompleteEntry { "find all files in a directory", "find -r *" },
		CompleteEntry { "write to a file", "cat"},
		CompleteEntry { "copy a directory", "cp -r"},
	}
}
