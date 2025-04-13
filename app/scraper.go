package app

import (
	"log"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

const file string = "./db.sql"

func Scrape(config common.Config) {

	dbConn := db.GetDbConnection(config["db_path"])
	defer dbConn.Close()

	linkRepository := db.NewLinkRepository(dbConn)
	pageRepository := db.NewPageRepository(dbConn)

	// if err := StoreLinks("2025-04-01", linkRepository); err != nil {
	// 	log.Fatalf("Error storing links: %v", err)
	// }

	if err := scrapeHtml("2025-04-01", linkRepository, pageRepository); err != nil {
		log.Fatalf("Error scraping HTML: %v", err)
	}

}
