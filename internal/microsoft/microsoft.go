package microsoft

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/tubopo/tick-tick-ticket/pkg/auth"
	"github.com/tubopo/tick-tick-ticket/pkg/calendar"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"github.com/tubopo/tick-tick-ticket/pkg/logger"
)

type Authenticator struct {
	Cfg *config.MicrosoftConfig
}

type authKey interface{}

func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	var authKey authKey

	if ctx.Value(authKey) != nil {
		return ctx, nil
	}

	apiUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(a.Cfg.TenantID))

	// Create form data with the required OAuth parameters.
	form := url.Values{}
	form.Add("client_id", a.Cfg.ClientID)
	form.Add("scope", "https://graph.microsoft.com/.default")
	form.Add("grant_type", "client_credentials")
	form.Add("client_secret", a.Cfg.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return ctx, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("could not authenticate: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, authKey, result.AccessToken)
	return ctx, nil
}

type Service struct {
	auth   auth.Authenticator
	Cfg    *config.MicrosoftConfig
	Logger logger.Logger
}

func NewService(cfg config.MicrosoftConfig, logger logger.Logger) calendar.Service {
	return &Service{
		auth:   &Authenticator{Cfg: &cfg},
		Cfg:    &cfg,
		Logger: logger,
	}
}
func (s *Service) GetCalendarEvents(start, end time.Time, ctx context.Context) ([]calendar.Event, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end must be set")
	}

	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/calendarView?startDateTime=%s&endDateTime=%s", startStr, endStr)
	s.Logger.Debug("Retrieving calendar events from %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	ctx, err = s.auth.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	var authKey authKey
	req.Header.Set("Authorization", "Bearer "+ctx.Value(authKey).(string))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []calendar.Event `json:"value"`
	}

	if resp.StatusCode == http.StatusOK {

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		s.Logger.Info("Got calendar events %d", len(result.Value))
	}

	return result.Value, nil
}

func (s *Service) CalcTotalDuration(events []calendar.Event) time.Duration {

	if len(events) == 0 {
		return 0
	}

	//sort events by start time
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	var totalDuration time.Duration
	var currentEnd time.Time

	for i, event := range events {
		if i == 0 || event.Start.After(currentEnd) {
			totalDuration += event.End.Sub(event.Start)
			currentEnd = event.End
		} else if event.End.After(currentEnd) {
			totalDuration += event.End.Sub(currentEnd)
			currentEnd = event.End
		}
		// If the event ends before the current end, it falls completely
		// within a previous event and does not add to the total duration.
	}

	return totalDuration
}
