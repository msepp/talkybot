package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

var selfRe *regexp.Regexp

func ircOnPrivMsg(event *irc.Event) {
	var channel string
	var msg string
	var shouldReact bool

	if len(event.Arguments) > 0 {
		channel = event.Arguments[0]
	}

	msg = event.Message()

	// check for mention, we react if mentioned. If not mentioned, give a small
	// chance of still reacting.
	shouldReact = strings.Contains(strings.ToLower(msg), strings.ToLower(app.cfg.Nick))
	if shouldReact == false {
		if rand.Intn(100) > 90 {
			shouldReact = true
		}
	}

	// Get conversation, create new one if doesn't exist if we should react.
	convo, ok := getConversation(channel, event.Nick, shouldReact)
	if !ok {
		return
	}

	// Strip obvious nick from start.
	msg = selfRe.ReplaceAllString(event.Message(), "")

	// feed new input into the conversation
	convo.NewInput(msg)
}

func ircOnBotReply(channel, nick, reply string, err error) {
	log.Printf("A reply to %s:%s, error: %v", channel, nick, err)

	if err != nil {
		// Terminate the conversation
		terminateConversation(channel, nick)
		return
	}

	var to = channel
	if to == "" {
		to = nick
	}

	if nick != "" {
		if rand.Intn(10) < 7 {
			reply = fmt.Sprintf("%s: %s", nick, reply)
		}
	}

	app.irc.Privmsg(to, reply)
}
