package logs

import "github.com/fatih/color"

var Quiet bool

func Info(message string) {
	if !Quiet {
		color.Cyan("[i] " + message)
	}
}

func Success(message string) {
	if !Quiet {
		color.Green("[+] " + message)
	}
}

func Error(message string) {
	if !Quiet {
		color.Red("[ERROR] " + message)
	}
}
