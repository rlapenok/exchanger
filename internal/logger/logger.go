package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rlapenok/exchanger/internal/config"
)

// init default logger for application
func init() {
	handler := slog.NewJSONHandler(
		os.Stdout,
		handlerOptions(slog.LevelDebug),
	)

	slog.SetDefault(slog.New(handler))
}

// redefine logger for application
func RedefineLogger(config config.LoggerConfig) {
	var (
		writer  io.Writer
		level   slog.Level
		handler slog.Handler
	)

	switch config.Output {
	case "stderr":
		writer = os.Stderr
	default:
		writer = os.Stdout
	}

	switch config.Level {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	switch config.Format {
	case "console":
		handler = slog.NewTextHandler(writer, handlerOptions(level))
	default:
		handler = slog.NewJSONHandler(writer, handlerOptions(level))
	}

	slog.SetDefault(slog.New(handler))
}

// handler options for logger
func handlerOptions(level slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, attr.Value.Time().Format(time.TimeOnly))
			}

			return attr
		},
	}
}
