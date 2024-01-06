package timeutil

import (
	"context"
	"time"
)

func SleepWithContext(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
