package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "./db.sql"

func main() {
	fmt.Print("hello !")
}

func getDbConnection(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
