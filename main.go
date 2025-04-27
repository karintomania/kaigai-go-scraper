package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/karintomania/kaigai-go-scraper/app"
	"github.com/karintomania/kaigai-go-scraper/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		mode        string
		date        string
		toStoreLink bool
		toScrape    bool
		toTranslate bool
		toGenerate  bool
		toPublish   bool
	)

	defaultDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	flag.StringVar(&mode, "mode", "scrape", "Mode of operation: migrate or scrape")
	flag.StringVar(&date, "date", defaultDate, "Date to scrape (YYYY-MM-DD)")
	flag.BoolVar(&toStoreLink, "l", false, "Flag to store links")
	flag.BoolVar(&toScrape, "s", false, "Flag to scrape")
	flag.BoolVar(&toTranslate, "t", false, "Flag to translate")
	flag.BoolVar(&toGenerate, "g", false, "Flag to generate")
	flag.BoolVar(&toPublish, "p", false, "Flag to generate")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		}))

	slog.SetDefault(logger)

	if mode == "migrate" {
		cmd.MigrateCmd()
	}

	if mode == "scrape" {
		app.Scrape(date, toStoreLink, toScrape, toTranslate, toGenerate, toPublish)
	}
}
