package db

import (
	"database/sql"
	"log"
)

const Rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

type TableCreater interface {
	CreateTable()
}

func GetDbConnection(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}

// Create testing database in tmp folder
// Don't forget to call deinit
func getTestEmptyDbConnection() (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatalln(err)
	}

	// TODO: remove cleanup
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func GetTestDbConnection() (*sql.DB, func()) {
	dbConn, cleanup := getTestEmptyDbConnection()

	Migrate(dbConn)

	return dbConn, cleanup
}

func Migrate(dbConn *sql.DB) {
	repos := []TableCreater{
		NewLinkRepository(dbConn),
		NewCommentRepository(dbConn),
		NewPageRepository(dbConn),
	}

	for _, repo := range repos {
		repo.CreateTable()
	}
}
