package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	if err := write(*c); err != nil {
		return err
	}
	return nil
}

func Read() (*Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return new(Config), err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return new(Config), err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	cfg := new(Config)
	if err = dec.Decode(cfg); err != nil {
		return new(Config), err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", home, configFilename), nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err = enc.Encode(cfg); err != nil {
		return err
	}
	return nil
}
