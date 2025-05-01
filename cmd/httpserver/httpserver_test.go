package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestHttpserver(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)

	date := "2025-10-15"

	pages := []db.Page{
		{
			Date:            "2025-10-15",
			Url:             "http://example.com/1",
			Title:           "Unpublished 1",
			TranslatedTitle: "未公開1",
			Published:       false,
		},
		{
			Date:            "2025-10-14",
			Url:             "http://example.com/2",
			Title:           "Unpublished 2",
			TranslatedTitle: "未公開2",
			Published:       false,
		},
		{
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

	mockPublishFunc := func() error {
		return nil
	}

	s := NewTestServer(dbConn, date, mockPublishFunc)

	common.MockEnv("server_port", "9999")
	rootUrl := fmt.Sprintf("http://localhost:%s/", common.GetEnv("server_port"))

	cli := &http.Client{}

	go func() {
		s.Start()
	}()

	// wait for the server to start
	time.Sleep(500 * time.Millisecond)

	t.Run("Get returns unpublished pages", func(t *testing.T) {
		req, err := http.NewRequest("GET", rootUrl, nil)
		require.NoError(t, err)

		response, err := cli.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)

		body := response.Body
		defer body.Close()
		htmlBytes, err := io.ReadAll(body)
		require.NoError(t, err)
		html := string(htmlBytes)

		t.Log(html)

		require.Contains(t, string(html), "未公開1")
		require.Contains(t, string(html), "未公開2")
		require.NotContains(t, string(html), "公開済1")
	})

	t.Run("Publish publishes", func(t *testing.T) {
		publishUrl := fmt.Sprintf("%spublish", rootUrl)

		req, err := http.NewRequest("POST", publishUrl, strings.NewReader("{}"))
		require.NoError(t, err)

		t.Logf("Before Do: Method=%s, URL=%s", req.Method, req.URL.String())

		response, err := cli.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)

		body := response.Body
		defer body.Close()
		htmlBytes, err := io.ReadAll(body)
		require.NoError(t, err)
		html := string(htmlBytes)

		t.Log(html)

		require.Contains(t, string(html), "Success")
	})
}
