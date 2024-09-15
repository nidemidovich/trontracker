package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Tronscan struct {
	APIKey string `envconfig:"TRONSCAN_API_KEY"`
}

func NewTronscan() (*Tronscan, error) {
	godotenv.Overload()

	tronscanConfig := &Tronscan{}

	err := envconfig.Process("TRONSCAN_API_KEY", tronscanConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get envs: %w", err)
	}

	return tronscanConfig, nil
}
