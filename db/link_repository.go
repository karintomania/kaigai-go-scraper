package db

import (
	"database/sql"
	"log/slog"
	// "fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Link struct {
	Id        int
	ExtId     string // ID of the original service e.g. HN
	Date      string
	URL       string
	Title     string
	Scraped   bool
	CreatedAt string
}

func (r *LinkRepository) Update(link *Link) {
	cmd := "UPDATE links SET ext_id = ?, date = ?, url = ?, title = ?, scraped = ? WHERE id = ?"
	_, err := r.db.Exec(cmd, link.ExtId, link.Date, link.URL, link.Title, link.Scraped, link.Id)
	if err != nil {
		log.Fatalln(err)
	}
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
		slog.Error("failed to insert link", "link", link)
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		slog.Error("failed to get last insert id link", "link", link)
		log.Fatalln(err)
	}

	link.Id = int(id)
}

func (r *LinkRepository) DoesExtIdExist(extId string) bool {
	query := "SELECT count(*) FROM links WHERE ext_id = ?"

	row := r.db.QueryRow(query, extId)

	result := 0

	err := row.Scan(&result)

	if err != nil {
		log.Fatalln(err)
	}

	return result > 0
}

func (r *LinkRepository) FindByDate(date string) []Link {
	query := "SELECT * FROM links WHERE date = ? AND scraped = 0"

	rows, err := r.db.Query(query, date)

	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	links := make([]Link, 0)

	for rows.Next() {
		var link Link

		err := rows.Scan(&link.Id, &link.ExtId, &link.Date, &link.URL, &link.Title, &link.Scraped, &link.CreatedAt)

		if err != nil {
			log.Fatalln(err)
		}

		links = append(links, link)
	}

	return links
}

func (r *LinkRepository) CreateTable() {
	cmd := `CREATE TABLE IF NOT EXISTS links(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	ext_id STRING NOT NULL,
	date STRING NOT NULL,
	url STRING NOT NULL,
	title STRING NOT NULL,
	scraped BOOLEAN NOT NULL DEFAULT 0,
	created_at STRING NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ'))
);
CREATE INDEX IF NOT EXISTS idx_links_date ON links(date);
CREATE INDEX IF NOT EXISTS idx_links_ext_id ON links(ext_id);
`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *LinkRepository) Drop() {
	cmd := "DROP TABLE IF EXISTS links"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *LinkRepository) Truncate() {
	cmd := "DELETE from links"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
