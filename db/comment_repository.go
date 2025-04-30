package db

import (
	"database/sql"
	// "fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	Id                int
	ExtCommentId      string
	PageId            int
	UserName          string
	Content           string
	TranslatedContent string
	Indent            int
	Reply             int
	Colour            string
	Score             int
	Translated        bool
	CommentedAt       string
	CreatedAt         string
}

func (r *CommentRepository) Update(comment *Comment) {
	cmd := `UPDATE comments SET ext_comment_id = ?, page_id = ?, user_name = ?, content = ?, translated_content = ?, indent = ?, reply = ?, colour = ?, score = ?, translated = ?, commented_at = ? WHERE id = ?`

	_, err := r.db.Exec(cmd,
		comment.ExtCommentId,
		comment.PageId,
		comment.UserName,
		comment.Content,
		comment.TranslatedContent,
		comment.Indent,
		comment.Reply,
		comment.Colour,
		comment.Score,
		comment.Translated,
		comment.CommentedAt,
		comment.Id)

	if err != nil {
		log.Fatalln(err)
	}
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Insert(comment *Comment) {
	cmd := `INSERT INTO comments (ext_comment_id, page_id, user_name, content, translated_content, indent, reply, colour, score, translated, commented_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(cmd,
		comment.ExtCommentId,
		comment.PageId,
		comment.UserName,
		comment.Content,
		comment.TranslatedContent,
		comment.Indent,
		comment.Reply,
		comment.Colour,
		comment.Score,
		comment.Translated,
		comment.CommentedAt)

	if err != nil {
		log.Fatalln(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatalln(err)
	}

	comment.Id = int(id)
}

func (r *CommentRepository) FindByPageId(pageId int) []Comment {
	query := "SELECT * FROM comments WHERE page_id = ?"

	rows, err := r.db.Query(query, pageId)

	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	comments := make([]Comment, 0)

	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.Id,
			&comment.ExtCommentId,
			&comment.PageId,
			&comment.UserName,
			&comment.Content,
			&comment.TranslatedContent,
			&comment.Indent,
			&comment.Reply,
			&comment.Colour,
			&comment.Score,
			&comment.Translated,
			&comment.CommentedAt,
			&comment.CreatedAt,
		)

		if err != nil {
			log.Fatalln(err)
		}

		comments = append(comments, comment)
	}

	return comments
}

func (r *CommentRepository) DoesExtIdExist(pageId int, extId string) bool {
	query := "SELECT count(*) FROM comments WHERE page_id = ? and ext_comment_id = ?"

	row := r.db.QueryRow(query, pageId, extId)

	count := 0

	err := row.Scan(&count)
	if err != nil {
		log.Fatalln(err)
	}

	return count > 0
}

func (r *CommentRepository) CreateTable() {
	cmd := `CREATE TABLE IF NOT EXISTS comments(
id INTEGER PRIMARY KEY AUTOINCREMENT,
ext_comment_id STRING NOT NULL,
page_id INTEGER NOT NULL,
user_name STRING NOT NULL,
content STRING NOT NULL,
translated_content STRING,
indent STRING,
reply STRING,
colour STRING,
score INTEGER,
translated BOOLEAN NOT NULL DEFAULT 0,
commented_at STRING,
created_at STRING NOT NULL DEFAULT (STRFTIME('%Y-%m-%dT%H:%M:%fZ'))
);
CREATE INDEX IF NOT EXISTS idx_comments_page_id ON comments(page_id);
CREATE INDEX IF NOT EXISTS idx_comments_ext_comment_id ON comments(ext_comment_id);
`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *CommentRepository) Drop() {
	cmd := "DROP TABLE IF EXISTS comments"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *CommentRepository) Truncate() {
	cmd := "DELETE from comments"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
