package app

import (
	"io"
	"strings"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestStoreLinks(t *testing.T) {
	dateStr := "2025-04-01"

	json := `[
{"id": "1", "points": 100, "link_text": "Link 1", "link": "http://link1.com", "dead": false},
{"id": "2", "points": 200, "link_text": "Link 2", "link": "http://link2.com", "dead": false},
{"id": "3", "points": 50, "link_text": "Link 3", "link": "http://link3.com", "dead": true},
{"id": "4", "points": 150, "link_text": "Link 4", "link": "http://link4.com", "dead": false},
{"id": "5", "points": 100, "link_text": "Link 5", "link": "http://link5.com", "dead": false}
]`

	mockCallHckrNews := func(dateStrPassed string) io.ReadCloser {
		require.Equal(t, dateStr, dateStrPassed)
		return io.NopCloser(strings.NewReader(json))
	}

	expected := []db.Link{
		{ExtId: "2", Date: dateStr, URL: "http://link2.com", Title: "Link 2", Scraped: false},
		{ExtId: "4", Date: dateStr, URL: "http://link4.com", Title: "Link 4", Scraped: false},
		{ExtId: "1", Date: dateStr, URL: "http://link1.com", Title: "Link 1", Scraped: false},
	}

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	common.MockEnv("top_link_num", "3")

	lr := db.NewLinkRepository(dbConn)

	storeLinks := NewTestStoreLinks(lr, mockCallHckrNews)
	err := storeLinks.run(dateStr)

	require.NoError(t, err)

	result := lr.FindByDate(dateStr)

	for i, gotLink := range result {
		require.Equal(t, expected[i].ExtId, gotLink.ExtId)
		require.Equal(t, expected[i].Date, gotLink.Date)
		require.Equal(t, expected[i].URL, gotLink.URL)
		require.Equal(t, expected[i].Title, gotLink.Title)
		require.Equal(t, expected[i].Scraped, gotLink.Scraped)
	}
}
