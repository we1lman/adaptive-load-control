package middleware

import (
	"time"

	"github.com/velma-andrei/adaptive-load-control/internal/limiter"
)

type statsObserver interface {
	Observe(d time.Duration)
}

func observe(l limiter.Limiter, d time.Duration) {
	if s, ok := l.(statsObserver); ok {
		s.Observe(d)
	}
}
