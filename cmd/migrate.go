package cmd

import (
	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

func Migrate(config common.Config) {
	dbConn := db.GetDbConnection(config["db_path"])
	defer dbConn.Close()

	lr := db.NewLinkRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)
	pr := db.NewPageRepository(dbConn)

	lr.CreateLinksTable()
	cr.CreateCommentsTable()
	pr.CreatePagesTable()
}
