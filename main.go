package main

import (
	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	migrate := false
	scrape := true

	config := map[string]string{
		"db_path": "./db.sql",
	}

	if migrate {
		cmd.Migrate(config)
	}

	if scrape {
		app.Scrape(config)
	}
}
