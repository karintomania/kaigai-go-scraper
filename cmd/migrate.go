package cmd

import (
	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

func Migrate() {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))
	defer dbConn.Close()

	lr := db.NewLinkRepository(dbConn)
	lr.CreateLinksTable()

	cr := db.NewCommentRepository(dbConn)
	cr.CreateCommentsTable()

	pr := db.NewPageRepository(dbConn)
	pr.CreatePagesTable()

}
