package db

import (
	"database/sql"
	"log"
)

func GetDbConnection(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
