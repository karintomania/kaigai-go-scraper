package app

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func TestTranslate(t *testing.T) {
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)

	pr.Insert(&db.Page{})
}
