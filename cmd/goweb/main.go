package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bajankristof/goweb/cli"
)

func main() {
	var err error

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errs := make(chan error, 1)
	go func() {
		errs <- cli.New().Run(ctx, os.Args)
	}()

	select {
	case <-ctx.Done():
		slog.DebugContext(ctx, "interrupted, bye now...")
		cancel()
		<-errs
		err = ctx.Err()
	case err = <-errs:
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		slog.ErrorContext(ctx, "fatal error", "err", err)
		os.Exit(1)
	}
}
