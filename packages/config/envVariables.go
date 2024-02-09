package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type EnvVariables struct {
	ConfigPath string `env:"CONFIG_PATH" envDefault:"dockerhook.yml"`
	Port       string `env:"PORT" envDefault:"8080"`
}

func LoadEnvVariables() (*EnvVariables, error) {
	fmt.Println("Loading env variables...")
	e := &EnvVariables{}

	if err := env.Parse(e); err != nil {
		fmt.Println("Something went wrong parsing the environment variables")
		return nil, err
	}

	fmt.Println("Environment variables loaded")

	return e, nil
}
