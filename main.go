package main

import (
	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	migrate := true
	scrape := false

	date := "2025-04-02"

	config := map[string]string{
		"db_path": "./db.sql",
	}

	if migrate {
		cmd.Migrate(config)
	}

	if scrape {
		app.Scrape(date, config)
	}
}
