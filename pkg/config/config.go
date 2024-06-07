package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Microsoft MicrosoftConfig `json:"microsoft"`
	Jira      JiraConfig      `json:"jira"`
}

type MicrosoftConfig struct {
	ApiKey string `json:"apiKey"`
}

type JiraConfig struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
	APIToken string `json:"apiToken"`
}

func Load(file string) (*Config, error) {
	cfg := &Config{}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
