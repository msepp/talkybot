package main

import (
	"log"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

func ircOnCommand(e *irc.Event) {
	parts := strings.SplitN(e.Message(), " ", 2)

	switch parts[0] {
	case "!quit":
		if len(parts) > 1 {
			app.irc.QuitMessage = strings.TrimPrefix(e.Message(), parts[1])
		}
		app.irc.Quit()

	case "!join":
		app.irc.Join(parts[1])

	case "!part":
		app.irc.Part(parts[1])

	default:
		log.Printf("Uknown command %v", parts)
	}
}
