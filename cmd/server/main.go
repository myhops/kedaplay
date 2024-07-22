package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"kedaplay/command"
	"kedaplay/service"
	"kedaplay/signalx"
)

func initLogger() *slog.Logger {
	h := slog.NewTextHandler(os.Stderr, nil)
	// h = slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

func main() {
	// Setup slog
	logger := initLogger()

	nctx, cancel := signalx.NotifyContext(context.Background())
	defer cancel()

	ctx := context.WithValue(nctx, service.SLoggerContextKey, logger)
	runner := command.NewServerCmd()
	runner.Run(ctx, os.Args, logger)

	s := signalx.CaughtSignal(nctx)
	if s != nil {
		log.Printf("caught signal: %s", s.String())
	}
}
