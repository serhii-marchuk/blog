package configs

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"os"
)

type Configs struct {
	WebPort  int `env:"WEB_PORT"`
	RestPort int `env:"REST_PORT"`
	Database DbConfig
}

func NewConfigs() *Configs {
	cfg := &Configs{}

	if godotenv.Load("./configs/.env") != nil {
		fmt.Println("Couldn't read .env file!")
		os.Exit(0)
	}

	if err := env.Parse(cfg); err != nil {
		fmt.Println("Couldn't parse .env file!")
		os.Exit(0)
	}

	return cfg
}
