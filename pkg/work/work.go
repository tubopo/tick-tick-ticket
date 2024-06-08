package work

import (
	"context"
	"time"
)

type Service interface {
	LogTime(jiraTicket string, duration time.Duration, ctx context.Context) error
}
