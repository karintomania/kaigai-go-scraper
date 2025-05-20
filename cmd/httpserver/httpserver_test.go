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

	errFlag := false
	mockPush := func() (string, error) {
		if errFlag {
			return "", fmt.Errorf("mock error")
		} else {
			return "Push succeed", nil
		}
	}

	s := NewTestServer(dbConn, date, mockPush)

	common.MockEnv("server_port", "9999")
	rootUrl := fmt.Sprintf("http://localhost:%s/", common.GetEnv("server_port"))
	common.MockEnv("server_host", "http://localhost")
	publishUrl := fmt.Sprintf("%spublish", rootUrl)

	cli := &http.Client{}

	go s.Start()

	// wait for the server to start
	time.Sleep(500 * time.Millisecond)

	t.Run("Get shows nothing to publish when no pages", func(t *testing.T) {
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

		require.Contains(t, string(html), "Nothing to publish")
	})

	t.Run("Get returns unpublished pages", func(t *testing.T) {
		for _, page := range pages {
			pr.Insert(&page)
		}

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

		require.Contains(t, string(html), "未公開1")
		require.Contains(t, string(html), "未公開2")
		require.NotContains(t, string(html), "公開済1")
	})

	t.Run("Publish publishes", func(t *testing.T) {
		req, err := http.NewRequest("POST", publishUrl, strings.NewReader("{}"))
		require.NoError(t, err)

		response, err := cli.Do(req)
		require.NoError(t, err)

		defer response.Body.Close()
		htmlBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		html := string(htmlBytes)

		require.Equal(t, http.StatusCreated, response.StatusCode)
		require.Contains(t, string(html), "Success")

		unpublished := pr.FindUnpublished()
		require.Len(t, unpublished, 0)
	})

	t.Run("Publish handles error", func(t *testing.T) {
		errFlag = true

		req, err := http.NewRequest("POST", publishUrl, strings.NewReader("{}"))
		require.NoError(t, err)

		response, err := cli.Do(req)
		require.NoError(t, err)

		defer response.Body.Close()
		htmlBytes, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		html := string(htmlBytes)

		require.Equal(t, http.StatusInternalServerError, response.StatusCode)
		require.Contains(t, string(html), "Something went wrong pushing git: mock error")
	})

	require.NoError(t, s.Shutdown())
}
