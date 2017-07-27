package main

import irc "github.com/thoj/go-ircevent"

func ircOnWelcome(e *irc.Event) {
	for _, ch := range app.cfg.Channels {
		app.irc.Join(ch)
	}
}
