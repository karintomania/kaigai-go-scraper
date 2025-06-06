package db

import (
	"database/sql"
	"log/slog"
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
	CreatedAt       string
}

func (r *PageRepository) Update(page *Page) error {
	cmd := `UPDATE pages SET ext_id = ?, date = ?, html = ?, title = ?, translated_title = ?, slug = ?, url = ?, ref_url = ?, tags = ?, translated = ?, published = ? WHERE id = ?`

	_, err := r.dbConn.Exec(cmd,
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

	return err
}

type PageRepository struct {
	dbConn *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository {
	return &PageRepository{dbConn: db}
}

func (r *PageRepository) Insert(page *Page) {
	cmd := `INSERT INTO pages (ext_id, date, html, title, translated_title, slug, url, ref_url, tags, translated, published) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.dbConn.Exec(cmd,
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
	)

	if err != nil {
		slog.Error("failed to store page", "page", page)
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatalln(err)
	}

	page.Id = int(id)
}

func (r *PageRepository) FindById(id int) *Page {
	query := "SELECT * FROM pages WHERE id = ?"

	row := r.dbConn.QueryRow(query, id)

	page := Page{}
	err := row.Scan(
		&page.Id,
		&page.ExtId,
		&page.Date,
		&page.Html,
		&page.Title,
		&page.TranslatedTitle,
		&page.Slug,
		&page.Url,
		&page.RefUrl,
		&page.Tags,
		&page.Translated,
		&page.Published,
		&page.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Panicln(err)
	}

	return &page
}

func (r *PageRepository) FindByDate(date string) []Page {
	query := "SELECT * FROM pages WHERE date = ?"

	rows, err := r.dbConn.Query(query, date)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	return r.scan(rows)

}

func (r *PageRepository) FindUntranslatedByDate(date string) []Page {
	query := "SELECT * FROM pages WHERE date = ? AND translated = 0"

	rows, err := r.dbConn.Query(query, date)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	return r.scan(rows)
}

func (r *PageRepository) FindUnpublished() []Page {
	query := "SELECT * FROM pages WHERE published = 0"

	rows, err := r.dbConn.Query(query)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	return r.scan(rows)
}

func (r *PageRepository) scan(rows *sql.Rows) []Page {
	pages := make([]Page, 0)

	for rows.Next() {
		var page Page

		err := rows.Scan(
			&page.Id,
			&page.ExtId,
			&page.Date,
			&page.Html,
			&page.Title,
			&page.TranslatedTitle,
			&page.Slug,
			&page.Url,
			&page.RefUrl,
			&page.Tags,
			&page.Translated,
			&page.Published,
			&page.CreatedAt,
		)

		if err != nil {
			log.Fatalln(err)
		}

		pages = append(pages, page)
	}

	return pages
}

func (r *PageRepository) CreateTable() {
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
	published BOOLEAN NOT NULL DEFAULT 0,
	created_at STRING NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ'))
);
CREATE INDEX idx_pages_date ON pages(date);`

	_, err := r.dbConn.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *PageRepository) Drop() {
	cmd := "DROP TABLE IF EXISTS pages"

	_, err := r.dbConn.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *PageRepository) Truncate() {
	cmd := "DELETE from pages"

	_, err := r.dbConn.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
