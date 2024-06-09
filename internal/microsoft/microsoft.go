package microsoft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/tubopo/tick-tick-ticket/pkg/auth"
	"github.com/tubopo/tick-tick-ticket/pkg/calendar"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"github.com/tubopo/tick-tick-ticket/pkg/logger"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Cfg      *config.MicrosoftConfig
	oauthCfg oauth2.Config
	state    string
}

type authTokenKey struct{}

var tokenCh = make(chan *oauth2.Token)

func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	a.state = "random-state"

	a.oauthCfg = oauth2.Config{
		ClientID:     a.Cfg.ClientID,
		ClientSecret: a.Cfg.ClientSecret,
		Scopes: []string{
			"https://graph.microsoft.com/Calendars.Read.Shared",
			"https://graph.microsoft.com/Calendars.ReadBasic",
			"https://graph.microsoft.com/Calendars.ReadWrite",
			"https://graph.microsoft.com/Calendars.ReadWrite.Shared",
			"https://graph.microsoft.com/User.Read"},
		RedirectURL: "http://localhost:8000/auth",
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", a.Cfg.TenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", a.Cfg.TenantID),
		},
	}

	http.HandleFunc("/auth", a.AuthCallBackHandler)
	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("Error starting server: %s\n", err.Error())
		}
	}()

	authURL := a.oauthCfg.AuthCodeURL(a.state, oauth2.AccessTypeOffline)
	fmt.Println("Please follow the URL to authenticate:", authURL)

	token := <-tokenCh
	var key authTokenKey
	ctx = context.WithValue(ctx, key, token.AccessToken)
	close(tokenCh)

	return ctx, nil
}

func (a *Authenticator) AuthCallBackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != a.state {
		http.Error(w, "Invalid state value", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := a.oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error exchanging token: %s\n", err.Error())
		return
	}

	tokenCh <- token

	w.Write([]byte("Authentication successful, you can close this window."))
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

type event struct {
	ID          string   `json:"id"`
	Subject     string   `json:"subject"`
	Start       startEnd `json:"start"`
	End         startEnd `json:"end"`
	Description string   `json:"description"`
}

type eventResponse struct {
	Value []event `json:"value"`
}

type startEnd struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

func (s *Service) GetCalendarEvents(start, end time.Time, ctx context.Context) ([]calendar.Event, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end must be set")
	}

	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/calendarView?startDateTime=%s&endDateTime=%s", url.QueryEscape(startStr), url.QueryEscape(endStr))
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

	var key authTokenKey
	req.Header.Set("Authorization", "Bearer "+ctx.Value(key).(string))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respText, _ := io.ReadAll(resp.Body)

		s.Logger.Debug("Got response:", resp.StatusCode, string(respText))

		return nil, errors.New("failed to retrieve calendar events")
	}

	var result = eventResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	s.Logger.Info("Got calendar events: ", len(result.Value))

	events := make([]calendar.Event, len(result.Value))

	const layout = "2006-01-02T15:04:05.0000000"
	for i, event := range result.Value {
		start, err := time.Parse(layout, event.Start.DateTime)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse(layout, event.End.DateTime)
		if err != nil {
			return nil, err
		}
		events[i] = calendar.Event{
			ID:    event.ID,
			Start: start,
			End:   end,
		}
	}

	return events, nil
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
