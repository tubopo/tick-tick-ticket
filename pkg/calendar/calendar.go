package calendar

import (
	"context"
	"time"
)

type Service interface {
	GetCalendarEvents(start, end time.Time, ctx context.Context) ([]Event, error)

	CalcTotalDuration(events []Event) time.Duration
}

type Event struct {
	ID    string
	Start time.Time
	End   time.Time
}
