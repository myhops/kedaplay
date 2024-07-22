package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"kedaplay/command"
	"kedaplay/service"
	"kedaplay/signalx"
)

func run(ctx context.Context, args []string, logger *slog.Logger) error {
	var cmd command.Cmd

	if len(args) < 2 {
		return errors.New("command missing")
	}
	switch args[1] {
	case "worker":
		cmd = command.NewWorkerCmd()
	case "server":
		cmd = command.NewServerCmd()
	default:
		return errors.New("bad command")
	}
	return cmd.Run(ctx, args, logger)
}

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
	if err := run(ctx, os.Args, logger); err != nil {
		log.Fatalf("run returned error: %s", err.Error())
	}

	s := signalx.CaughtSignal(nctx)
	if s != nil {
		log.Printf("caught signal: %s", s.String())
	}
}
