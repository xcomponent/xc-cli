package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os/user"
)

type Config struct {
	ServerUrl string
}

func LoadConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	var conf Config
	filePath := fmt.Sprintf("%s/%s", usr.HomeDir, SettingsFile)
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
