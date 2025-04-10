package app

import (
	"log"

	"github.com/karintomania/kaigai-go-scraper/db"
)

const file string = "./db.sql"

func Scrape() {

	dbConn := db.GetDbConnection(file)
	defer dbConn.Close()

	linkRepository := db.NewLinkRepository(dbConn)

	if err := StoreLinks("2025-04-01", linkRepository); err != nil {
		log.Fatalf("Error storing links: %v", err)
	}
}
