package main

import (
	"encoding/json"
	"errors"
	"os"
)

// Config describes target system and recipients
type Config struct {
	// CleverBotAPIKey is the API key for cleverbot API
	CleverBotAPIKey string
	// Nick is the nickname of the bot
	Nick string
	// Realname is the realname of the bot
	Realname string
	// User is the username of the bot
	Username string
	// Channels contains a list of channels to join
	Channels []string
	// Server is the server to connect to
	Server string
	// OwnerNick is the nick of the bots owner
	OwnerNick string
	// OwnerUser is the user of the bots owner
	OwnerUser string
	// EngageRandomly tells if the bot should enage in conversations randomly
	EngageRandomly bool
}

func readConfig(filename string) (cfg *Config, err error) {
	var f *os.File

	if f, err = os.Open(filename); err != nil {
		return
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&cfg); err != nil {
		return
	}

	if err = validateConfig(cfg); err != nil {
		return
	}

	return
}

func validateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if cfg.CleverBotAPIKey == "" {
		return errors.New("CleverBotAPIKey not set")
	}

	if cfg.Server == "" {
		return errors.New("Server not set")
	}

	if cfg.Nick == "" {
		return errors.New("Nick not set")
	}

	if cfg.Username == "" {
		return errors.New("Username not set")
	}

	return nil
}
