package config

import (
	"errors"
	"os"
)

// Config holds runtime settings sourced from the environment.
type Config struct {
	BotToken string
}

// Load reads configuration from environment variables.
// BOT_TOKEN is required; an empty value is a fatal misconfiguration.
func Load() (Config, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return Config{}, errors.New("BOT_TOKEN is required")
	}
	return Config{BotToken: token}, nil
}
