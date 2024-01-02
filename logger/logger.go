package logger

import (
	"context"
	"log/slog"
	"os"
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
