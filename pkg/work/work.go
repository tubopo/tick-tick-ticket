package work

import (
	"context"
	"time"
)

type Service interface {
	LogTime(duration time.Duration, date time.Time, ctx context.Context) error
}
