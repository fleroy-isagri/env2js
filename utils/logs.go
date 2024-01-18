package utils

import (
	"github.com/fatih/color"
)

func LogError(title string, log string) {
	color.New(color.Bold).Add(color.FgRed).Println(title + " : " + log)
}

func LogSuccess(title string, log string) {
	color.New(color.Bold).Print(title)
	color.Green(log)
}
