package middleware

import (
	"context"
	"time"

	"github.com/velma-andrei/adaptive-load-control/internal/limiter"
)

type Middleware func(ctx context.Context, priority int, handle func() error) error

func New(l limiter.Limiter) Middleware {
	return func(ctx context.Context, priority int, handle func() error) error {
		release, decision := l.Acquire(ctx, priority)
		if decision == limiter.Reject {
			return ErrOverloaded
		}
		start := time.Now()
		err := handle()
		observe(l, time.Since(start))
		release()
		return err
	}
}

var ErrOverloaded = errorString("resource exhausted")

type errorString string

func (e errorString) Error() string { return string(e) }
