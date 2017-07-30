package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	irc "github.com/thoj/go-ircevent"
)

var useDebug bool
var app = struct {
	cfg           *Config
	irc           *irc.Connection
	conversations map[string]*Conversation
	mutex         sync.Mutex
}{
	conversations: map[string]*Conversation{},
}

func terminateConversation(channel, nick string) {
	key := nick + channel

	app.mutex.Lock()
	defer app.mutex.Unlock()

	log.Printf("Terminating conversation with key %s", key)

	if _, ok := app.conversations[key]; ok {
		delete(app.conversations, key)
	}
}

func getConversation(channel, nick string, createIfNotExist bool) (*Conversation, bool) {
	var convo *Conversation
	var ok bool
	var key string

	key = nick + channel

	app.mutex.Lock()
	defer app.mutex.Unlock()

	if convo, ok = app.conversations[key]; !ok {
		log.Printf("convo with key %s doesn't exist, should create? %t", key, createIfNotExist)
		if createIfNotExist {
			convo = NewConversation(channel, nick, ircOnBotReply)
			app.conversations[key] = convo
		} else {
			return nil, false
		}
	}

	return convo, true
}

func init() {
	var configPath string
	flag.BoolVar(&useDebug, "debug", false, "Enable debug messages")
	flag.StringVar(&configPath, "config", "config.json", "Path to configuration")
	flag.Parse()

	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Loading configuration failed: %s", err)
	}

	app.cfg = config

	// compile self patter
	selfRe = regexp.MustCompile(`(?i)\s?` + app.cfg.Nick + `[:,>]\s?`)
}

func main() {
	app.irc = irc.IRC(app.cfg.Nick, app.cfg.Username)

	if useDebug {
		app.irc.VerboseCallbackHandler = true
		app.irc.Debug = true
	}

	app.irc.AddCallback("001", ircOnWelcome)
	app.irc.AddCallback("PRIVMSG", ircOnPrivMsg)

	err := app.irc.Connect(app.cfg.Server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	// Start a timer for reaping old conversations
	go func() {
		var tc = time.Tick(time.Second)
		for {
			select {
			case <-tc:
				for _, c := range app.conversations {
					if c.Idle() > (time.Minute * 2) {
						log.Printf("Idle conversation, terminating")
						terminateConversation(c.Channel(), c.Nick())
					}
				}
			}
		}
	}()
	app.irc.Loop()

	log.Printf("Exiting. Bye!")
}
