package main

import (
	"flag"
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

func Test_parseFlags_Version(t *testing.T) {
	// Arrange
	version = "1.0.0"
	expectedOutput := "version : 1.0.0\ncommit  : \ndate    : \nbuiltBy : \n"

	// Act
	config, output, _ := ParseFlags("prog", []string{"-version"})

	// Assert
	if config.version != true {
		t.Errorf("config.version should be true")
	}
	if output != expectedOutput {
		t.Errorf("output = %v, expected %v", output, expectedOutput)
	}
}

func Test_parseFlags_Help(t *testing.T) {
	// Arrange
	version = "1.0.0"
	expectedOutput := "Usage of prog:\n  -version\n    \tDisplay version and exit\n"

	// Act
	config, output, err := ParseFlags("prog", []string{"-help"})

	// Assert
	if config != nil {
		t.Errorf("config should be nil")
	}
	if err == nil {
		t.Errorf("err should not be nil")
	}
	if err != flag.ErrHelp {
		t.Errorf("err should be errorFlagHelp")
	}
	if output != expectedOutput {
		t.Errorf("output = %v, expected %v", output, expectedOutput)
	}
}
