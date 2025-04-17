package app

import (
	"fmt"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func TestTranslatePage(t *testing.T) {
	page := &db.Page{
		Title: "test title",
	}

	comments := make([]db.Comment, 0, 10)

	for i := 0; i < 10; i++ {
		comment := db.Comment{
			Content: fmt.Sprintf("test comment %d", i),
		}

		comments = append(comments, comment)
	}

	translatePage(page, comments)
}
