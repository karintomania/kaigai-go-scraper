package app

import (
	"log"
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
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
		if err := StoreLinks(date, linkRepository); err != nil {
			log.Panicf("Error storing links: %v", err)
		}
	}

	if toScrape {
		if err := scrapeHtml(
			date,
			linkRepository,
			pageRepository,
			commentRepository); err != nil {
			log.Panicf("Error scraping HTML: %v", err)
		}
	}

	if toTranslate {
		if err := translate(
			date,
			pageRepository,
			commentRepository); err != nil {
			log.Panicf("Error translating: %v", err)
		}
	}

	if toGenerate {
		ag := NewArticleGenerator(pageRepository, commentRepository)

		if err := ag.generateArticles(date); err != nil {
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
