package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/velma-andrei/adaptive-load-control/internal/limiter"
	"github.com/velma-andrei/adaptive-load-control/internal/middleware"
)

func main() {
	addr := flag.String("addr", "localhost:9090", "listen address")
	flag.Parse()

	l := limiter.New(limiter.Config{
		InitialLimit: 16,
		MinLimit:     1,
		MaxLimit:     128,
	})
	mw := middleware.New(l)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("listening on %s", *addr)

	simulate := func() error {
		return mw(ctx, 5, func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(20 * time.Millisecond):
				return nil
			}
		})
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := simulate(); err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("request failed: %v", err)
			}
		}
	}
}
