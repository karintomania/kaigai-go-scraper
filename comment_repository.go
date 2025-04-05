package main

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
	Indent            string
	Reply             string
	Colour            string
	Score             int
}

func (r *CommentRepository) Update(comment Comment) {
	cmd := `UPDATE Comments SET ext_comment_id = ?, page_id = ?, user_name = ?, content = ?, translated_content = ?, indent = ?, reply = ?, colour = ?, score = ? WHERE id = ?`

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

func (r *CommentRepository) Insert(comment Comment) {
	cmd := `INSERT INTO Comments (ext_comment_id, page_id, user_name, content, translated_content, indent, reply, colour, score) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(cmd,
		comment.ExtCommentId,
		comment.PageId,
		comment.UserName,
		comment.Content,
		comment.TranslatedContent,
		comment.Indent,
		comment.Reply,
		comment.Colour,
		comment.Score)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *CommentRepository) FindByPageId(pageId int) []Comment {
	query := "SELECT * FROM Comments WHERE page_id = ?"

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
			&comment.Score)

		if err != nil {
			log.Fatalln(err)
		}

		comments = append(comments, comment)
	}

	return comments
}

func (r *CommentRepository) CreateCommentsTable() {
	cmd := `CREATE TABLE IF NOT EXISTS Comments(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
ext_comment_id STRING NOT NULL,
page_id INTEGER NOT NULL,
user_name STRING NOT NULL,
content STRING NOT NULL,
translated_content STRING,
indent STRING,
reply STRING,
colour STRING,
score INTEGER
	)`

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}

func (r *CommentRepository) DropCommentsTable() {
	cmd := "DROP TABLE IF EXISTS Comments"

	_, err := r.db.Exec(cmd)

	if err != nil {
		log.Fatalln(err)
	}
}
