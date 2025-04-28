package scrape

import (
	"fmt"
	"strings"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestTranslateTitle(t *testing.T) {
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

	translatedTitle := "翻訳後タイトル"
	translatedTags := []string{"タグ１", "タグ２"}

	callAIMock := func(prompt string) string {
		tagsStr := strings.Join(translatedTags, `","`)

		return fmt.Sprintf(`{"title":"%s","tags":["%s"]}`, translatedTitle, tagsStr)
	}

	translateTitle(page, comments, callAIMock)

	if want, got := translatedTitle, page.TranslatedTitle; want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}

	gotTags := strings.Split(page.Tags, ",")

	for i, tag := range translatedTags {
		if want, got := tag, gotTags[i]; want != got {
			t.Errorf("Expected %s, got %s", want, got)
		}
	}
}

func TestTranslateComment(t *testing.T) {
	comments := make([]db.Comment, 0, 10)
	for i := 0; i < 5; i++ {
		comment := db.Comment{
			Id:      i + 1,
			Content: fmt.Sprintf("test comment %d", i+1),
		}
		comments = append(comments, comment)
	}

	callAIMock := func(prompt string) string {
		return `{"comments":[
			{"id": 1, "content": "翻訳コメント 1", "score": 10},
			{"id": 2, "content": "翻訳コメント 2", "score": 20},
			{"id": 3, "content": "翻訳コメント 3", "score": 30},
			{"id": 4, "content": "翻訳コメント 4", "score": 40},
			{"id": 5, "content": "翻訳コメント 5", "score": 50}
		]}`
	}

	_, err := translateCommentChunk("翻訳タイトル", comments, callAIMock)

	require.NoError(t, err)

	for i, comment := range comments {
		if fmt.Sprintf("翻訳コメント %d", i+1) != comment.TranslatedContent ||
			(i+1)*10 != comment.Score ||
			true != comment.Translated {
			t.Errorf("Unexpected comment: %v", comment)
		}
	}
}

func TestSanitizeTranslatedComment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<test>", "＜test＞"}, // escape <>
		// escape quotes
		{`'quoted', "double-quoted"`, `’quoted’, ”double-quoted”`},
		// keep br
		{"<br>", "<br>"},
		// convert new line to br
		{"\n", "<br>"},
	}
	for _, test := range tests {
		result := sanitizeTranslatedComment(test.input)
		if result != test.expected {
			t.Errorf("sanitizeTranslatedComment(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
