package main

import (
	"strings"

	irc "github.com/thoj/go-ircevent"
)

func ircOnWelcome(e *irc.Event) {
	for _, ch := range app.cfg.Channels {
		app.irc.Join(ch)
	}
}

func ircOnPrivMsg(e *irc.Event) {
	if e.Nick == app.cfg.OwnerNick && e.User == app.cfg.OwnerUser {
		if strings.HasPrefix(e.Message(), "!") {
			ircOnCommand(e)
			return
		}
	}

	ircOnConversationMessage(e)
}
