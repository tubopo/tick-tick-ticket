package microsoft

import (
	"context"
	"time"

	"github.com/tubopo/tick-tick-ticket/pkg/auth"
	"github.com/tubopo/tick-tick-ticket/pkg/calendar"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"github.com/tubopo/tick-tick-ticket/pkg/logger"
)

type Authenticator struct {
	Cfg *config.MicrosoftConfig
}

func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	//todo: implement this
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
	//todo: implement this
	var events []calendar.Event

	return events, nil
}
