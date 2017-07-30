package main

import (
	"log"
	"sync"
	"time"

	cleverbot "github.com/ugjka/cleverbot-go"
)

type ConversationCallback func(channel, nick, reply string, err error)

// Conversation is a single conversation.
type Conversation struct {
	// channel that the message should be posted to
	channel string
	// nick the bot is currently talking with.
	nick string
	// cb is the CleverBot API session handle
	bot *cleverbot.Session
	// nextInput is the next input we're gonna fetch a response for.
	nextInput string
	// mutex for safe handling
	mutex sync.RWMutex
	// flag telling that response is being fetched
	gettingReply bool
	// callback to call when a reply is generated.
	cb ConversationCallback
	// lastActive tells when the convo was last active
	lastActive time.Time
}

// NewConversation initializes and returns a ready-to-use Conversation
func NewConversation(channel, nick string, callback ConversationCallback) *Conversation {
	return &Conversation{
		nick:    nick,
		channel: channel,
		bot:     cleverbot.New(app.cfg.CleverBotAPIKey),
		cb:      callback,
	}
}

// NewInput feeds new input into the conversation buffer
func (c *Conversation) NewInput(line string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.nextInput += line
	c.lastActive = time.Now()

	if c.gettingReply == false {
		c.gettingReply = true
		go c.GetReply(c.nextInput)
		c.nextInput = ""
	}
}

// GetReply retrieves a reply to given string from cleverbot
func (c *Conversation) GetReply(to string) {
	answer, err := c.bot.Ask(to)
	log.Printf("Got answer: %s, delay will be %dms", answer, len(answer)*200)
	if err == nil {
		time.Sleep(time.Duration((len(answer) * 200)) * time.Millisecond)
	}
	c.cb(c.channel, c.nick, answer, err)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lastActive = time.Now()
	c.gettingReply = false
}

// Nick returns conversation nick
func (c *Conversation) Nick() string {
	return c.nick
}

// Channel returns conversation channel
func (c *Conversation) Channel() string {
	return c.channel
}

// Idle returns how long the conversation has been idle
func (c *Conversation) Idle() time.Duration {
	n := time.Now()
	return n.Sub(c.lastActive)
}
