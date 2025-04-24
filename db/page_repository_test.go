package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageRepository(t *testing.T) {
	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	pr := NewPageRepository(db)

	pr.CreateTable()

	t.Run("Insert", func(t *testing.T) {
		pr.Truncate()

		page := Page{
			ExtId:           "12345",
			Date:            "2025-01-01",
			Html:            "<html></html>",
			Title:           "Test Title",
			TranslatedTitle: "Translated Test Title",
			Slug:            "test-title",
			Url:             "https://example.com/12345",
			RefUrl:          "https://example.com",
			Tags:            "tag1,tag2",
			Translated:      true,
			Published:       true,
		}

		pr.Insert(&page)

		pages := pr.FindByDate("2025-01-01")

		require.Equal(t, 1, len(pages), "Expected %d page, got %d", 1, len(pages))

		created := pages[0]

		require.Equal(t, page.ExtId, created.ExtId)
		require.Equal(t, page.Date, created.Date)
		require.Equal(t, page.Html, created.Html)
		require.Equal(t, page.Title, created.Title)
		require.Equal(t, page.TranslatedTitle, created.TranslatedTitle)
		require.Equal(t, page.Slug, created.Slug)
		require.Equal(t, page.Url, created.Url)
		require.Equal(t, page.RefUrl, created.RefUrl)
		require.Equal(t, page.Tags, created.Tags)
		require.Equal(t, page.Translated, created.Translated)
		require.Equal(t, page.Published, created.Published)
	})

	t.Run("Update", func(t *testing.T) {
		pr.Truncate()
		// Insert initial page
		page := Page{
			ExtId:           "12345",
			Date:            "2025-01-01",
			Html:            "<html></html>",
			Title:           "Test Title",
			TranslatedTitle: "Translated Test Title",
			Slug:            "test-title",
			Url:             "https://example.com/12345",
			RefUrl:          "https://example.com",
			Tags:            "tag1,tag2",
		}
		pr.Insert(&page)

		// Update the page
		updatedPage := &Page{
			Id:              page.Id,
			ExtId:           "54321",
			Date:            "2025-01-02",
			Html:            "<html>Updated</html>",
			Title:           "Updated Title",
			TranslatedTitle: "Updated Translated Title",
			Slug:            "updated-title",
			Url:             "https://example.com/54321",
			RefUrl:          "https://example-updated.com",
			Tags:            "tag3,tag4",
			Translated:      true,
			Published:       true,
		}
		pr.Update(updatedPage)

		// Verify the update
		pages := pr.FindByDate("2025-01-02")
		require.Equal(t, 1, len(pages), "Expected %d page, got %d", 1, len(pages))

		updated := pages[0]
		require.Equal(t, updatedPage.ExtId, updated.ExtId)
		require.Equal(t, updatedPage.Date, updated.Date)
		require.Equal(t, updatedPage.Html, updated.Html)
		require.Equal(t, updatedPage.Title, updated.Title)
		require.Equal(t, updatedPage.TranslatedTitle, updated.TranslatedTitle)
		require.Equal(t, updatedPage.Slug, updated.Slug)
		require.Equal(t, updatedPage.Url, updated.Url)
		require.Equal(t, updatedPage.RefUrl, updated.RefUrl)
		require.Equal(t, updatedPage.Tags, updated.Tags)
		require.Equal(t, updatedPage.Translated, updated.Translated)
		require.Equal(t, updatedPage.Published, updated.Published)
	})

	t.Run("FindByDate returns nothing for no record", func(t *testing.T) {
		pr.Truncate()

		pages := pr.FindByDate("2025-01-01")

		require.Equal(t, 0, len(pages), "Expected %d page, got %d", 0, len(pages))
	})

}
