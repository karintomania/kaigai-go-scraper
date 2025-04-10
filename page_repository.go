package main

import (
	"database/sql"
	// "fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Id              int
	ExtId           string
	Date            string
	Html            string
	Title           string
	TranslatedTitle string
	Slug            string
	Url             string
	RefUrl          string
	Tags            string // comma separated tags
	Translated      bool
	Published       bool
}

func (r *PageRepository) Update(page Page) {
	cmd := `UPDATE pages SET ext_id = ?, date = ?, html = ?, title = ?, translated_title = ?, slug = ?, url = ?, ref_url = ?, tags = ?, translated = ?, published = ? WHERE id = ?`

	_, err := r.db.Exec(cmd,
		page.ExtId,
		page.Date,
		page.Html,
		page.Title,
		page.TranslatedTitle,
		page.Slug,
		page.Url,
		page.RefUrl,
		page.Tags,
		page.Translated,
		page.Published,
		page.Id)

	if err != nil {
		log.Fatalln(err)
	}
}

type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{db: db}
}

func (r *PageRepository) Insert(page *Page) {
	cmd := `INSERT INTO pages (ext_id, date, html, title, translated_title, slug, url, ref_url, tags, translated, published) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(cmd,
		page.ExtId,
		page.Date,
		page.Html,
		page.Title,
		page.TranslatedTitle,
		page.Slug,
		page.Url,
		page.RefUrl,
		page.Tags,
		page.Translated,
		page.Published)

	if err != nil {
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatalln(err)
	}

	page.Id = int(id)
}

func (r *PageRepository) FindByDate(date string) []Page {
	query := "SELECT * FROM pages WHERE date = ?"

	rows, err := r.db.Query(query, date)

	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	Pages := make([]Page, 0)

	for rows.Next() {
		var Page Page

		err := rows.Scan(
			&Page.Id,
			&Page.ExtId,
			&Page.Date,
			&Page.Html,
			&Page.Title,
			&Page.TranslatedTitle,
			&Page.Slug,
			&Page.Url,
			&Page.RefUrl,
			&Page.Tags,
			&Page.Translated,
			&Page.Published)

		if err != nil {
			log.Fatalln(err)
		}

		Pages = append(Pages, Page)
	}

	return Pages

}

func (r *PageRepository) CreatePagesTable() {
	cmd := `CREATE TABLE IF NOT EXISTS pages(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ext_id STRING NOT NULL,
			date STRING NOT NULL,
			html STRING NOT NULL,
			title STRING,
			translated_title STRING,
			slug STRING,
			url STRING NOT NULL UNIQUE,
			ref_url STRING,
			tags STRING,
translated BOOLEAN NOT NULL DEFAULT 0,
published BOOLEAN NOT NULL DEFAULT 0
	)`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *PageRepository) DropPagesTable() {
	cmd := "DROP TABLE IF EXISTS pages"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
