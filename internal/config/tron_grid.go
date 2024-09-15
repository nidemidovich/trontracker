package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type TronGrid struct {
	APIKey string `envconfig:"TRON_GRID_API_KEY"`
}

func NewTronGrid() (*TronGrid, error) {
	godotenv.Overload()

	tronGridConfig := &TronGrid{}

	err := envconfig.Process("TRON_GRID_API_KEY", tronGridConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get envs: %w", err)
	}

	return tronGridConfig, nil
}
