package scrape

import (
	"log"
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

func Scrape(
	date string,
	toStoreLink bool,
	toScrape bool,
	toTranslate bool,
	toGenerate bool,
	toPublish bool,
) {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))
	defer dbConn.Close()

	linkRepository := db.NewLinkRepository(dbConn)
	pageRepository := db.NewPageRepository(dbConn)
	commentRepository := db.NewCommentRepository(dbConn)

	if toStoreLink {
		sl := NewStoreLinks(linkRepository)

		if err := sl.run(date); err != nil {
			log.Panicf("Error storing links: %v", err)
		}
	}

	if toScrape {
		scrapeHtml := NewScrapeHtml(linkRepository, pageRepository, commentRepository)

		if err := scrapeHtml.run(date); err != nil {
			log.Panicf("Error scraping HTML: %v", err)
		}
	}

	if toTranslate {
		tp := NewTranslatePage(pageRepository, commentRepository, external.CallGemini)

		if err := tp.run(date); err != nil {
			log.Panicf("Error translating: %v", err)
		}
	}

	if toGenerate {
		ag := NewGenerateArticle(pageRepository, commentRepository)

		if err := ag.run(date); err != nil {
			log.Panicf("Error translating: %v", err)
		}
	}

	if toPublish {
		p := NewPublish(pageRepository)

		if err := p.run(date); err != nil {
			log.Panicf("Error publishing: %v", err)
		}
	}

	slog.Info("Scraping completed")
}
