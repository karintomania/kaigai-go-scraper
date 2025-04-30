package scrape

import (
	"fmt"
	"strings"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestTraslatePage(t *testing.T) {

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)

	t.Run("TranslateTitle translates the title", func(t *testing.T) {
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

		callAiMock := func(prompt string) (string, error) {
			tagsStr := strings.Join(translatedTags, `","`)

			return fmt.Sprintf(`{"title":"%s","tags":["%s"]}`, translatedTitle, tagsStr), nil
		}

		tr := NewTranslatePage(pr, cr, callAiMock)

		err := tr.translateTitle(page, comments)

		require.NoError(t, err)

		require.Equal(t, translatedTitle, page.TranslatedTitle)

		gotTags := strings.Split(page.Tags, ",")

		for i, tag := range translatedTags {
			require.Equal(t, tag, gotTags[i])
		}
	})

	t.Run("TranslateCommentChunk translates comments", func(t *testing.T) {

		comments := make([]db.Comment, 0, 10)
		for i := 0; i < 5; i++ {
			comment := db.Comment{
				Id:      i + 1,
				Content: fmt.Sprintf("test comment %d", i+1),
			}
			comments = append(comments, comment)
		}

		callAiMock := func(prompt string) (string, error) {
			return `{"comments":[
			{"id": 1, "content": "翻訳コメント 1", "score": 10},
			{"id": 2, "content": "翻訳コメント 2", "score": 20},
			{"id": 3, "content": "翻訳コメント 3", "score": 30},
			{"id": 4, "content": "翻訳コメント 4", "score": 40},
			{"id": 5, "content": "翻訳コメント 5", "score": 50}
		]}`, nil
		}

		tr := NewTranslatePage(pr, cr, callAiMock)
		_, err := tr.translateCommentChunk("翻訳タイトル", comments)

		require.NoError(t, err)

		for i, comment := range comments {
			require.Equal(t, fmt.Sprintf("翻訳コメント %d", i+1), comment.TranslatedContent)
			require.Equal(t, (i+1)*10, comment.Score)
			require.Equal(t, true, comment.Translated)
		}
	})

	t.Run("retry retries correctly", func(t *testing.T) {
		t.Run("retry just once when successful", func(t *testing.T) {
			tr := NewTranslatePage(pr, cr, nil)

			// Mock the function to be retried
			called := 0
			mockFunc := func() error {
				called++
				return nil
			}

			// Call the retry function
			err := tr.retry(mockFunc, 3)

			// Check the result and error
			require.NoError(t, err)
			require.Equal(t, 1, called)
		})

		t.Run("retry correct times on error", func(t *testing.T) {
			tr := NewTranslatePage(pr, cr, nil)

			called := 0
			mockErrorFunc := func() error {
				called++
				return fmt.Errorf("try: %d", called)
			}

			// Call the retry function
			err := tr.retry(mockErrorFunc, 3)

			// Check the result and error
			require.EqualError(t, err, "failed to run function after tried 3 times")
			require.Equal(t, 3, called)
		})
	})
}
