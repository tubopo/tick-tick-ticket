package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tubopo/tick-tick-ticket/pkg/auth"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"github.com/tubopo/tick-tick-ticket/pkg/logger"
	"github.com/tubopo/tick-tick-ticket/pkg/work"
)

type Authenticator struct {
	Cfg *config.JiraConfig
}

type authKey struct{}

type workLogPayload struct {
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}

func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	var authKey authKey

	if ctx.Value(authKey) != nil {
		return ctx, nil
	}

	if a.Cfg.APIToken == "" {
		return nil, errors.New("API token is not set")
	}

	return context.WithValue(ctx, authKey, a.Cfg.APIToken), nil
}

type Service struct {
	auth       auth.Authenticator
	Cfg        *config.JiraConfig
	jiraTicket string
	Logger     logger.Logger
}

func NewService(cfg config.JiraConfig, jiraTicket string, log logger.Logger) work.Service {
	return &Service{
		auth:       &Authenticator{Cfg: &cfg},
		Cfg:        &cfg,
		jiraTicket: jiraTicket,
		Logger:     log,
	}
}

func (s *Service) LogTime(duration time.Duration, ctx context.Context) error {
	if s.jiraTicket == "" {
		return errors.New("jira ticket is not set")
	}

	if duration == 0 {
		return errors.New("duration is not set")
	}

	url := fmt.Sprintf("https://%s/rest/api/2/issue/%s/worklog", s.Cfg.Domain, s.jiraTicket)
	s.Logger.Debug("Logging time to %s", url)

	payload, err := json.Marshal(workLogPayload{TimeSpentSeconds: int(duration.Seconds())})
	if err != nil {
		return err
	}

	s.Logger.Debug("Payload: %+v", payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, err = s.auth.Authenticate(ctx)
	if err != nil {
		return err
	}

	var authKey authKey
	req.Header.Set("Authorization", "Bearer "+ctx.Value(authKey).(string))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		s.Logger.Info("Logged time to %s  - %v", url, resp.StatusCode)
		return nil
	}

	s.Logger.Debug("Got response: %v", resp.StatusCode)

	return errors.New("failed to log time to JIRA")
}
