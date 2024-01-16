package main

import "testing"

func Test_DefineFilePath_MultiFileChoice(t *testing.T) {
	want := "tests/example-1.js"
	if got, _ := DefineFilePath("./tests", "example"); got != want {
		t.Errorf("DefineFilePath() = %v, want %v", got, want)
	}
}

func Test_DefineFilePath_OneFileChoice(t *testing.T) {
	want := "tests/example-2.js"
	if got, _ := DefineFilePath("./tests", "example-2"); got != want {
		t.Errorf("DefineFilePath() = %v, want %v", got, want)
	}
}

func Test_DefineFilePath_ErrorNoFileFound(t *testing.T) {
	if _, err := DefineFilePath("./tests", "toto"); err.Error() != "No file found with pattern: tests/toto*.js" {
		t.Errorf("DefineFilePath() does not return an error : " + err.Error())
	}
}
