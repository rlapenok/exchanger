package main

import (
	"log/slog"

	"github.com/rlapenok/exchanger/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("failed to run application", "error", err)
	}
}
