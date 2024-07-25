package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"syscall"

	"kedaplay/command"
	"kedaplay/signalx"
)

type options struct {
}

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
	return cmd.Run(ctx, args[1:], logger)
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

	nctx, cancel := signalx.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ctx := nctx
	if err := run(ctx, os.Args, logger); err != nil {
		log.Fatalf("run returned error: %s", err.Error())
	}

	s := signalx.CaughtSignal(nctx)
	if s != nil {
		logger.Info("caught signal", slog.String("signal", s.String()))
	} else {
		logger.Info("caught signal", slog.String("detail", "signal is nil"))
	}
	log.Print("stopped")
}
