package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Microsoft MicrosoftConfig `json:"calendar"`
	Jira      JiraConfig      `json:"jira"`
}

type MicrosoftConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	TenantID     string `json:"tenantId"`
}

type JiraConfig struct {
	Domain   string `json:"domain"`
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
