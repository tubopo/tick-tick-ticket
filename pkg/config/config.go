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
	TenantID     string   `json:"tenantId"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Scopes       []string `json:"scopes"`
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
