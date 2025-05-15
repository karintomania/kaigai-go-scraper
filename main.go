package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/karintomania/kaigai-go-scraper/cmd"
	"github.com/karintomania/kaigai-go-scraper/cmd/httpserver"
	"github.com/karintomania/kaigai-go-scraper/cmd/scrape"
	"github.com/karintomania/kaigai-go-scraper/cmd/tweets"
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
		toRunAll    bool
	)

	defaultDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	flag.StringVar(&mode, "mode", "scrape", "Mode of operation: migrate/scrape/server")
	flag.StringVar(&date, "date", defaultDate, "Date to scrape (YYYY-MM-DD)")
	flag.BoolVar(&toRunAll, "run-all", false, "Alias for -l -s -t -g -p")
	flag.BoolVar(&toStoreLink, "l", false, "Flag to store links")
	flag.BoolVar(&toScrape, "s", false, "Flag to scrape")
	flag.BoolVar(&toTranslate, "t", false, "Flag to translate")
	flag.BoolVar(&toGenerate, "g", false, "Flag to generate")
	flag.BoolVar(&toPublish, "p", false, "Flag to generate")
	flag.Parse()

	if toRunAll {
		toStoreLink = true
		toScrape = true
		toTranslate = true
		toGenerate = true
		toPublish = true
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		}))

	slog.SetDefault(logger)

	if mode == "migrate" {
		cmd.MigrateCmd()
	}

	if mode == "scrape" {
		scrape.Scrape(date, toStoreLink, toScrape, toTranslate, toGenerate, toPublish)
	}

	if mode == "server" {
		s := httpserver.NewServer()

		s.Start()
	}

	if mode == "tweet" {
		tweetCmd := tweets.NewPostScheduledCmd()

		if err := tweetCmd.Run(date); err != nil {
			slog.Error("Error posting scheduled tweets", "error", err)
		}
	}
}
