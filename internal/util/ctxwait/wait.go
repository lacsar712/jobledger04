package ctxwait

import (
	"context"
	"time"
)

func Until(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	time.Sleep(d)
	return ctx.Err()
}
