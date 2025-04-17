package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/karintomania/kaigai-go-scraper/common"
)

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
	file, err := os.CreateTemp("", "testing.sql")

	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}

	common.MockEnv("db_path", file.Name())

	db, err := sql.Open("sqlite3", file.Name())

	if err != nil {
		log.Fatalln(err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(file.Name())
	}

	return db, cleanup
}

func GetTestDbConnection() (*sql.DB, func()) {
	db, cleanup := getTestEmptyDbConnection()

	Migrate()

	return db, cleanup
}

func Migrate() {
	dbConn := GetDbConnection(common.GetEnv("db_path"))
	defer dbConn.Close()

	repos := []TableCreater{
		NewLinkRepository(dbConn),
		NewCommentRepository(dbConn),
		NewPageRepository(dbConn),
	}

	for _, repo := range repos {
		repo.CreateTable()
	}
}
