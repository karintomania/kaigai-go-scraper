package db

import (
	"database/sql"
	"log"
)

type Tweet struct {
	Id          int
	PageId      int    // comma separated tags
	Date        string // ISO8601 (YYYY-MM-DD)
	Content     string
	ScheduledAt string // ISO8601 (YYYY-MM-DD)
	Published   bool   // default false
	CreatedAt   string
}

type TweetRepository struct {
	db *sql.DB
}

func NewTweetRepository(db *sql.DB) *TweetRepository {
	return &TweetRepository{db: db}
}

func (r *TweetRepository) FindById(id int) *Tweet {
	query := "SELECT * FROM tweets WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var tweet Tweet
	err := row.Scan(
		&tweet.Id,
		&tweet.PageId,
		&tweet.Date,
		&tweet.Content,
		&tweet.ScheduledAt,
		&tweet.Published,
		&tweet.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Panicln(err)
	}

	return &tweet
}

func (r *TweetRepository) Insert(tweet *Tweet) {
	cmd := `INSERT INTO tweets (page_id, date, content, scheduled_at, published) VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(cmd,
		tweet.PageId,
		tweet.Date,
		tweet.Content,
		tweet.ScheduledAt,
		tweet.Published,
	)

	if err != nil {
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatalln(err)
	}

	tweet.Id = int(id)
}

func (r *TweetRepository) Update(tweet *Tweet) {
	cmd := `UPDATE tweets SET page_id = ?, date = ?, content = ?, scheduled_at = ?, published = ? WHERE id = ?`

	_, err := r.db.Exec(cmd,
		tweet.PageId,
		tweet.Date,
		tweet.Content,
		tweet.ScheduledAt,
		tweet.Published,
		tweet.Id,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *TweetRepository) CreateTable() {
	cmd := `CREATE TABLE IF NOT EXISTS tweets(
id INTEGER PRIMARY KEY AUTOINCREMENT,
page_id INTEGER NOT NULL,
date STRING NOT NULL,
content STRING NOT NULL,
scheduled_at STRING NOT NULL,
published BOOLEAN NOT NULL DEFAULT 0,
created_at STRING NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ'))
);
CREATE INDEX IF NOT EXISTS idx_scheduled_at ON tweets(scheduled_at)`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

// Find Tweets which is scheduled before the given date
// and not published
func (r *TweetRepository) FindUnpublishedByScheduledDate(dateBy string) []Tweet {
	query := "SELECT * FROM tweets WHERE scheduled_at <= ? AND published = 0 order by id"

	rows, err := r.db.Query(query, dateBy)

	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	tweets := make([]Tweet, 0)

	for rows.Next() {
		var tweet Tweet

		err := rows.Scan(
			&tweet.Id,
			&tweet.PageId,
			&tweet.Date,
			&tweet.Content,
			&tweet.ScheduledAt,
			&tweet.Published,
			&tweet.CreatedAt,
		)

		if err != nil {
			log.Fatalln(err)
		}

		tweets = append(tweets, tweet)
	}

	return tweets
}

func (r *TweetRepository) Drop() {
	cmd := "DROP TABLE IF EXISTS tweets"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *TweetRepository) Truncate() {
	cmd := "DELETE from tweets"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
