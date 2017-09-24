package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os/user"
	"os"
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
		switch err.(type) {
		default:
			return nil, err
		case *os.PathError:
			return &Config{}, nil
		}
	}

	return &conf, nil
}
