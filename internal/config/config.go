package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil{
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	err := writeConfig(*c)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgFilePath, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	} else {
		return homeDir + "/" + configFileName, nil
	}
}
