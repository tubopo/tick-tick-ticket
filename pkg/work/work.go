package work

import (
	"context"
	"time"
)

type Service interface {
	LogTime(duration time.Duration, ctx context.Context) error
}
