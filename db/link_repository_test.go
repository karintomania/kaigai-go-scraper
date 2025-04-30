package db

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/stretchr/testify/require"
)

func TestLinkRepository(t *testing.T) {
	common.SetLogger()

	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	lr := NewLinkRepository(db)
	lr.CreateTable()

	t.Run("Insert", func(t *testing.T) {
		lr.Truncate()

		link := Link{
			ExtId:   "12345",
			Date:    "2025-01-01",
			URL:     "https://example.com/12345",
			Title:   "Example Title",
			Scraped: false,
		}

		lr.Insert(&link)

		links := lr.FindByDate("2025-01-01")

		require.Equal(t, 1, len(links), "Expected 1 link")

		created := links[0]

		require.Equal(t, link.ExtId, created.ExtId, "ExtId mismatch")
		require.Equal(t, link.Date, created.Date, "Date mismatch")
		require.Equal(t, link.URL, created.URL, "URL mismatch")
		require.Equal(t, link.Title, created.Title, "Title mismatch")
		require.Equal(t, link.Scraped, created.Scraped, "Scraped mismatch")
	})

	t.Run("FindByIdReturnNothing", func(t *testing.T) {
		lr.Truncate()

		links := lr.FindByDate("2025-01-01")

		require.Equal(t, 0, len(links), "Expected 0 links")
	})

	t.Run("DoesExtIdExist", func(t *testing.T) {
		lr.Truncate()

		link := Link{
			ExtId:   "1",
			Date:    "2025-01-01",
			URL:     "https://example.com/12345",
			Title:   "Example Title",
			Scraped: false,
		}

		lr.Insert(&link)

		require.True(t, lr.DoesExtIdExist("1"))
		require.False(t, lr.DoesExtIdExist("2"))
	})
}
