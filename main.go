package main

import (
	"flag"
	"fmt"
	"log"

	irc "github.com/thoj/go-ircevent"
	cleverbot "github.com/ugjka/cleverbot-go"
)

var app = struct {
	cfg       *Config
	cleverbot *cleverbot.Session
	irc       *irc.Connection
}{}

func init() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.json", "Path to configuration")
	flag.Parse()

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Loading configuration failed: %s", err)
	}

	app.cfg = config
}

func main() {
	app.irc = irc.IRC(app.cfg.Nick, app.cfg.Username)
	app.irc.VerboseCallbackHandler = true
	app.irc.Debug = true
	app.irc.AddCallback("001", ircOnWelcome)
	app.irc.AddCallback("JOIN", func(e *irc.Event) {
		log.Printf("join:\n%+v", e)
		go func() {
			app.cleverbot = cleverbot.New(app.cfg.CleverBotAPIKey)
			answer, _ := app.cleverbot.Ask("Hi, How are you?")
			app.irc.Privmsg("#sweetiechan", answer)
		}()
	})

	err := app.irc.Connect(app.cfg.Server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	app.irc.Loop()
}
