package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bajankristof/watchbowl/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errs := make(chan error, 1)
	go func() {
		errs <- cli.New().Run(ctx, os.Args)
	}()

	var err error
	select {
	case <-ctx.Done():
		stop()
		slog.DebugContext(ctx, "interrupted, bye now...")
		<-errs
		err = ctx.Err()
	case err = <-errs:
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
}
