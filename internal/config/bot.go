package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Bot struct {
	Token string `envconfig:"BOT_API_TOKEN"`
}

func NewBot() (*Bot, error) {
	godotenv.Overload()

	botConfig := &Bot{}

	err := envconfig.Process("TOKEN", botConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get envs: %w", err)
	}

	return botConfig, nil
}
