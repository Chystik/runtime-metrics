package main

import (
	"os"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func parseEnv(cfg *config.AgentConfig) error {
	var envFile string
	var err error
	// load ENV from file if ENVIRONMENT val is set and parses it values

	if osEnv := os.Getenv("ENVIRONMENT"); osEnv != "" {
		switch osEnv {
		case "dev":
			envFile = "agent_dev.env"
		case "prod":
			envFile = "agent_prod.env"
		case "stage":
			envFile = "agent_stage.env"
		}

		// load ENV from *.env file
		err = godotenv.Load(envFile)
		if err != nil {
			panic(err)
		}
	}

	return env.Parse(cfg)
}
