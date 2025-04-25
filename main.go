package main

import (
	"log/slog"
	"os"

	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	migrate := false
	scrape := true

	date := "2025-04-01"

	logger := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		}))

	slog.SetDefault(logger)

	if migrate {
		cmd.MigrateCmd()
	}

	if scrape {
		app.Scrape(date)
	}
}
