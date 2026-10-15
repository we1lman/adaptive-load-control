package limiter

import (
	"context"
	"testing"
	"time"
)

func TestAcquireWithinLimit(t *testing.T) {
	l := New(Config{InitialLimit: 2})
	_, d := l.Acquire(context.Background(), PriorityHigh)
	if d != Accept {
		t.Fatal("expected accept")
	}
}

func TestAcquireRejectsWhenInfeasible(t *testing.T) {
	l := New(Config{InitialLimit: 4})
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	_, d := l.Acquire(ctx, PriorityHigh)
	if d != Reject {
		t.Fatal("expected reject")
	}
}

func TestStatsMean(t *testing.T) {
	s := &Stats{}
	s.Observe(10 * time.Millisecond)
	s.Observe(30 * time.Millisecond)
	if s.Mean() != 20*time.Millisecond {
		t.Fatal("unexpected mean")
	}
}
