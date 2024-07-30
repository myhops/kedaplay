package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/myhops/kedaplay/command"
	"github.com/myhops/kedaplay/signalx"
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

	ctx := nctx
	runner := command.NewWorkerCmd(&command.WorkerConfig{})
	runner.Run(ctx)

	s := signalx.CaughtSignal(nctx)
	if s != nil {
		logger.Info("caught signal", slog.String("signal", s.String()))
	}
}
