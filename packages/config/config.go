package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigFile struct {
	Config config
	Auth   auth
}

type config struct {
	Enable        bool   `yaml:"enable"`
	LabelBased    bool   `yaml:"labelBased"`
	DefaultAction string `yaml:"defaultAction" validate:"oneof=pull start stop restart"`
}

type auth struct {
	Tokens     []string `yaml:"tokens"`
	TokensFile string   `yaml:"tokensFile"`
	Groups     []string `yaml:"groups"`
}

func LoadConfig(configPath string) (*ConfigFile, error) {
	var c = &ConfigFile{}
	fmt.Println("Reading the configuration file...")
	configYml, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Println("Something went wrong reading the configuration file, please check if the path is correct")
		return nil, err
	}

	err = yaml.Unmarshal(configYml, c)

	if err != nil {
		fmt.Println("Something went wrong parsing the configuration file, please check if the file is correct")
		return nil, err
	}

	validate := validator.New()

	err = validate.Struct(c)

	if err != nil {
		fmt.Println("Something went wrong validating the configuration file, please check if the file is correct")
		return nil, err
	}

	fmt.Println("Configuration file read successfully!")

	return c, nil
}
