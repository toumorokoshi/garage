package garage

import "testing"

func TestGetCompleteEntryFromString(t *testing.T) {
	output := getCompleteEntryFromString("gfind: foo; bar")
	desiredOutput := CompleteEntry{ "foo", "bar" }
	if output.Message != desiredOutput.Message {
		t.Logf("Error! expected %s, got %s", desiredOutput, output)
		t.Fail()
	}
	if output.Command != desiredOutput.Command {
		t.Logf("Error! expected %s, got %s", desiredOutput, output)
		t.Fail()
	}
}
