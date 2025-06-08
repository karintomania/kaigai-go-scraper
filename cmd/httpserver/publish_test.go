package httpserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)
	// tr := db.NewTweetRepository(dbConn)

	pages := []db.Page{
		{
			Id:              1,
			Date:            "2025-10-15",
			Url:             "http://example.com/1",
			Title:           "Unpublished 1",
			TranslatedTitle: "未公開1",
			Published:       false,
		},
		{
			Id:              2,
			Date:            "2025-10-14",
			Url:             "http://example.com/2",
			Title:           "Unpublished 2",
			TranslatedTitle: "未公開2",
			Published:       false,
		},
		{
			Id:              3,
			Date:            "2025-10-13",
			Url:             "http://example.com/3",
			Title:           "Published 1",
			TranslatedTitle: "公開済1",
			Published:       true,
		},
	}

	for _, page := range pages {
		pr.Insert(&page)
	}

	errFlag := false
	mockPush := func() (string, error) {
		if errFlag {
			return "", fmt.Errorf("mock error")
		} else {
			return "Push succeed", nil
		}
	}

	mockScheduleTweet := func(dateStr string, pageIds []int) error {
		today := time.Now().Format("2006-01-02")

		require.Equal(t, today, dateStr)
		require.Equal(t, 2, len(pageIds))
		for i, id := range pageIds {
			t.Logf("test")
			require.Equal(t, i+1, id)
		}

		return nil
	}

	handler := PublishHandler{
		push:     mockPush,
		schedule: mockScheduleTweet,
		pr:       pr,
	}

	formData := url.Values{}

	formData.Add("page_ids", "1")
	formData.Add("page_ids", "2")

	r := httptest.NewRequest(http.MethodPost, "/publish", strings.NewReader(formData.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.handle(w, r)

	res := w.Result()

	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)
}
