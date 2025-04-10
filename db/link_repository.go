package db

import (
	"database/sql"
	// "fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Link struct {
	Id      int
	ExtId   string // ID of the original service e.g. HN
	Date    string
	URL     string
	Title   string
	Scraped bool
}

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// Persist Link and add the ID
func (r *LinkRepository) Insert(link *Link) {
	cmd := "INSERT INTO links (ext_id, date, url, title, scraped) VALUES (?, ?, ?, ?, ?)"

	result, err := r.db.Exec(cmd, link.ExtId, link.Date, link.URL, link.Title, link.Scraped)
	if err != nil {
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatalln(err)
	}

	link.Id = int(id)
}

func (r *LinkRepository) FindByDate(date string) []Link {
	query := "SELECT * FROM links WHERE date = ?"

	rows, err := r.db.Query(query, date)

	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	links := make([]Link, 0)

	for rows.Next() {
		var link Link

		err := rows.Scan(&link.Id, &link.ExtId, &link.Date, &link.URL, &link.Title, &link.Scraped)

		if err != nil {
			log.Fatalln(err)
		}

		links = append(links, link)
	}

	return links

}

func (r *LinkRepository) CreateLinksTable() {
	cmd := `CREATE TABLE IF NOT EXISTS links(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ext_id STRING NOT NULL,
			date STRING NOT NULL,
			url STRING NOT NULL UNIQUE,
title STRING NOT NULL,
scraped BOOLEAN NOT NULL DEFAULT 0
	)`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *LinkRepository) DropLinksTable() {
	cmd := "DROP TABLE IF EXISTS links"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
