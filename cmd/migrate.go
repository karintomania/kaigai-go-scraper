package cmd

import (
	"github.com/karintomania/kaigai-go-scraper/db"
)

func MigrateCmd() {
	db.Migrate()
}
