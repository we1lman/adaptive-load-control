package limiter

import (
	"context"
	"sync"
	"time"
)

type Decision int

const (
	Accept Decision = iota
	Reject
)

const (
	PriorityLow  = 0
	PriorityHigh = 9
)

type Limiter interface {
	Acquire(ctx context.Context, priority int) (release func(), decision Decision)
}

type Config struct {
	InitialLimit int
	MinLimit     int
	MaxLimit     int
}

type Stats struct {
	mu        sync.Mutex
	count     int64
	totalTime time.Duration
}

func (s *Stats) Observe(d time.Duration) {
	s.mu.Lock()
	s.count++
	s.totalTime += d
	s.mu.Unlock()
}

func (s *Stats) Mean() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 {
		return 0
	}
	return s.totalTime / time.Duration(s.count)
}

type combinedLimiter struct {
	mu       sync.Mutex
	inFlight int
	cfg      Config
	stats    *Stats
}

func New(cfg Config) Limiter {
	if cfg.InitialLimit <= 0 {
		cfg.InitialLimit = 8
	}
	if cfg.MinLimit <= 0 {
		cfg.MinLimit = 1
	}
	if cfg.MaxLimit < cfg.InitialLimit {
		cfg.MaxLimit = cfg.InitialLimit
	}
	return &combinedLimiter{cfg: cfg, stats: &Stats{}}
}

func (l *combinedLimiter) limit() int {
	return l.cfg.InitialLimit
}

func (l *combinedLimiter) feasible(ctx context.Context) bool {
	deadline, ok := ctx.Deadline()
	if !ok {
		return true
	}
	remaining := time.Until(deadline)
	return remaining > 3*l.stats.Mean()
}

func (l *combinedLimiter) Acquire(ctx context.Context, priority int) (func(), Decision) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inFlight >= l.limit() && priority <= PriorityLow {
		return func() {}, Reject
	}
	if !l.feasible(ctx) {
		return func() {}, Reject
	}
	l.inFlight++
	release := func() {
		l.mu.Lock()
		l.inFlight--
		l.mu.Unlock()
	}
	return release, Accept
}
