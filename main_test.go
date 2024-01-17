package main

import (
	"path/filepath"
	"testing"
)

func Test_DefineFilePath_MultiFileChoice(t *testing.T) {
	want := "tests/example-1.js"
	if got, _ := DefineFilePath("./tests", "example"); filepath.ToSlash(got) != want {
		t.Errorf("DefineFilePath() = %v, want %v", got, want)
	}
}

func Test_DefineFilePath_OneFileChoice(t *testing.T) {
	want := "tests/example-2.js"
	if got, _ := DefineFilePath("./tests", "example-2"); filepath.ToSlash(got) != want {
		t.Errorf("DefineFilePath() = %v, want %v", got, want)
	}
}

func Test_DefineFilePath_ErrorNoFileFound(t *testing.T) {
	if _, err := DefineFilePath("./tests", "toto"); filepath.ToSlash(err.Error()) != "No file found with pattern: "+filepath.ToSlash("tests/toto*.js") {
		t.Errorf("DefineFilePath() does not return an error : " + filepath.ToSlash(err.Error()))
	}
}
