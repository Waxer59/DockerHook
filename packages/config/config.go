package config

type action int

const (
	Pull action = iota
	Start
	Stop
	Restart
)

type Config struct {
	LabelBased    bool
	DefaultAction action
	Auth          auth
}

type auth struct {
	GroupTokens []groupToken
	AuthGroups  []authGroup
	TokensFile  string
}

type groupToken struct {
	Name       string
	Tokens     []string
	FileTokens []string
}

type authGroup struct {
	Name                                                               string
	HavePullAccess, HaveStartAccess, HaveStopAccess, HaveRestartAccess bool
}

func (c *Config) LoadConfig(configPath string) error {
	return nil
}
