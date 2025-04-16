package main

import (
	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	migrate := false
	scrape := false

	date := "2025-04-02"

	if migrate {
		cmd.Migrate()
	}

	if scrape {
		app.Scrape(date)
	}
}
