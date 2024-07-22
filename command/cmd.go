package command

import (
	"context"
	"log/slog"
)

type Cmd interface {
	Run(ctx context.Context, args []string, logger *slog.Logger) error
}
