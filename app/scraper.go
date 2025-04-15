package app

import (
	"log"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

const file string = "./db.sql"

func Scrape(date string, config common.Config) {
	dbConn := db.GetDbConnection(config["db_path"])
	defer dbConn.Close()

	linkRepository := db.NewLinkRepository(dbConn)
	pageRepository := db.NewPageRepository(dbConn)
	commentRepository := db.NewCommentRepository(dbConn)

	// if err := StoreLinks(date, linkRepository); err != nil {
	// 	log.Fatalf("Error storing links: %v", err)
	// }

	if err := scrapeHtml(
		date,
		linkRepository,
		pageRepository,
		commentRepository); err != nil {
		log.Fatalf("Error scraping HTML: %v", err)
	}
}
