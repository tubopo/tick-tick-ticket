package auth

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
)

type Authenticator interface {
	Authenticate(ctx context.Context) (context.Context, error)
}

type TokenStore struct {
	Token *oauth2.Token
}

func (s *TokenStore) Save() error {
	if s.Token == nil {
		return errors.New("token is nil")
	}

	os.UserHomeDir()

	data, err := json.Marshal(s.Token)
	if err != nil {
		return err
	}
	return os.WriteFile(getTokenFilePath(), data, 0600)
}

func (s *TokenStore) Load() error {
	data, err := os.ReadFile(getTokenFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &s.Token)
}

func (s *TokenStore) IsValid() bool {
	return s.Token != nil && s.Token.Expiry.After(time.Now())
}

func getTokenFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	appDir := filepath.Join(homeDir, ".tick-tick-ticket")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		os.MkdirAll(appDir, 0700)
	}

	return filepath.Join(appDir, "token.json")
}
