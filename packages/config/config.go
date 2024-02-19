package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var actions = []string{"start", "stop", "restart", "pull"}

type ConfigFile struct {
	Config config
	Auth   auth
}

type config struct {
	RemoveOldImage bool   `yaml:"removeOldImage"`
	LabelBased     bool   `yaml:"labelBased"`
	DefaultAction  string `yaml:"defaultAction" validate:"oneof=pull start stop restart"`
}

type auth struct {
	Enable     bool     `yaml:"enable"`
	Tokens     []string `yaml:"tokens"`
	TokensFile string   `yaml:"tokensFile"`
	Groups     []string `yaml:"groups"`
}

func LoadConfig(configPath string) (*ConfigFile, error) {
	c := &ConfigFile{ // default config
		Config: config{
			RemoveOldImage: false,
			LabelBased:     false,
			DefaultAction:  actions[3],
		},
		Auth: auth{
			Enable:     false,
			Tokens:     []string{},
			Groups:     []string{},
			TokensFile: "",
		},
	}
	fmt.Println("Looking for configuration file")
	configYml, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Println("No configuration file was found")
		return c, nil
	}

	err = yaml.Unmarshal(configYml, c)

	if err != nil {
		fmt.Println("Something went wrong parsing the configuration file, please check if the file is correct")
		return nil, err
	}

	fileTokens, err := c.loadFileTokens()

	if err != nil {
		fmt.Println("Something went wrong reading the tokens file, please check if the path is correct")
		return nil, err
	}

	c.Auth.Tokens = append(c.Auth.Tokens, fileTokens...)

	validate := validator.New()

	err = validate.Struct(c)

	if err != nil {
		fmt.Println("Something went wrong validating the configuration file, please check if the file is correct")
		return nil, err
	}

	fmt.Println("Configuration file read successfully!")

	return c, nil
}

func (c *ConfigFile) GetGroups() map[string][]string {
	groups := make(map[string][]string)

	for _, group := range c.Auth.Groups {
		groupName := strings.Split(group, ":")[0] // groupName:<action1>,<action2>
		actionsRaw := strings.Split(group, ":")[1]
		actions := strings.Split(actionsRaw, ",") // [action1, action2]
		groups[groupName] = actions
	}

	return groups
}

func (c *ConfigFile) GetTokenActions(token string) []string {
	groups := c.GetGroups()
	groupName := ""

	for _, tokenSearch := range c.Auth.Tokens {
		if strings.Contains(tokenSearch, token) {
			groupName = strings.Split(tokenSearch, ":")[0]
			break
		}
	}

	// If no group is specified in the token, the token will have full access
	if groupName == "" {
		return actions
	}

	return groups[groupName]
}

func (c *ConfigFile) loadFileTokens() ([]string, error) {
	var tokens []string
	if c.Auth.TokensFile != "" {
		fmt.Println("Reading tokens from file")
		file, err := os.ReadFile(c.Auth.TokensFile)

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, strings.Split(string(file), "\n")...)
	}

	return tokens, nil
}
