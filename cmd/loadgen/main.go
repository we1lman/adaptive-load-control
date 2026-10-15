package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rate := flag.Float64("rate", 100, "requests per second")
	duration := flag.Duration("duration", 30*time.Second, "load duration")
	target := flag.String("target", "localhost:9090", "target address")
	flag.Parse()

	log.Printf("target=%s rate=%v duration=%v", *target, *rate, *duration)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
}
