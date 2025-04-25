package app

import (
	"log"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

func Scrape(date string) {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))
	defer dbConn.Close()

	linkRepository := db.NewLinkRepository(dbConn)
	pageRepository := db.NewPageRepository(dbConn)
	commentRepository := db.NewCommentRepository(dbConn)

	toStoreLink := false
	toScrape := true
	toTranslate := true
	toGenerate := true

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
}
