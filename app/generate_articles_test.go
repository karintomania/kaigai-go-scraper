package app

import (
	"fmt"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestGenerateArticle(t *testing.T) {
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)

	mockGetImage := func() string {
		return "mock_image.jpg"
	}

	mockGetColour := func() string {
		return "#000000"
	}

	common.MockEnv("minimum_colour_score", "20")

	page := &db.Page{
		Id:              1,
		TranslatedTitle: "Test Title",
		Slug:            "test-title",
		Url:             "http://example.com",
		RefUrl:          "http://example.com/ref",
		Tags:            "tag1,tag2",
	}

	comments := make([]db.Comment, 0)
	for i := 1; i < 28; i++ {
		c := db.Comment{
			Id:                i,
			UserName:          fmt.Sprintf("User-%d", i),
			TranslatedContent: fmt.Sprintf("コメント %d", i),
			Score:             i,
			CommentedAt:       fmt.Sprintf("2025-12-%02dT00:01:02.345Z", i),
		}

		comments = append(comments, c)
	}

	ag := NewTestArticleGenerator(
		pr,
		cr,
		mockGetImage,
		mockGetColour,
	)

	// Test case 1: Basic functionality
	article, err := ag.generateArticle("2025-01-01", page, comments)

	require.NoError(t, err)

	require.Contains(t, article, "date = '2025-01-01T00:00:00'")
	require.Contains(t, article, "months = '2025/01'")
	require.Contains(t, article, "title = 'Test Title'")
	require.Contains(t, article, `tags = ["tag1", "tag2"]`)
	require.Contains(t, article, `featureimage = 'mock_image.jpg'`)
	require.Contains(t, article, `引用元：[http://example.com](http://example.com)`)

	for i := 1; i < 28; i++ {
		date := fmt.Sprintf("2025/12/%02d 00:01:02", i)
		colour := ""
		if i >= 20 {
			colour = "#000000"
		}

		require.Contains(t, article, fmt.Sprintf(`{{<matomeQuote body="コメント %d" userName="User-%d" createdAt="%s" color="%s">}}`, i, i, date, colour))
	}

	require.Contains(t, article, `{{< details summary="もっとコメントを表示（1）">}}`)
	require.Contains(t, article, `{{< details summary="もっとコメントを表示（2）">}}`)
	require.NotContains(t, article, `{{< details summary="もっとコメントを表示（3）">}}`)
}
