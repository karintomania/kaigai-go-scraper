package main

import (
	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Migrate("./db.sql")

	app.Scrape()
}
