package logger

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

type loggerKey struct{}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := ctx.Value(loggerKey{})
	if logger == nil {
		return newLogger()
	}
	return logger.(*slog.Logger)
}
