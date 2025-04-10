package cmd

import "github.com/karintomania/kaigai-go-scraper/db"

func Migrate(file string) {
	dbConn := db.GetDbConnection(file)
	defer dbConn.Close()

	lr := db.NewLinkRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)
	pr := db.NewPageRepository(dbConn)

	lr.CreateLinksTable()
	cr.CreateCommentsTable()
	pr.CreatePagesTable()
}
