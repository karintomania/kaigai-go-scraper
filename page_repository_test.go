package main

import (
	"os"
	"testing"
)

func TestPageRepositoryInsert(t *testing.T) {
	file, err := os.CreateTemp("", "PageRepository.sql")

	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	db := getDbConnection(file.Name())

	pr := NewPageRepository(db)

	pr.CreatePagesTable()

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

	if want, got := 1, len(pages); want != got {
		t.Fatalf("Expected %d page, got %d", want, got)
	}

	created := pages[0]

	if created.ExtId != page.ExtId ||
		created.Date != page.Date ||
		created.Html != page.Html ||
		created.Title != page.Title ||
		created.TranslatedTitle != page.TranslatedTitle ||
		created.Slug != page.Slug ||
		created.Url != page.Url ||
		created.RefUrl != page.RefUrl ||
		created.Tags != page.Tags ||
		created.Translated != page.Translated ||
		created.Published != page.Published {
		t.Fatalf("Expected %v, got %v", page, created)
	}
}

func TestPageRepositoryUpdate(t *testing.T) {
	file, err := os.CreateTemp("", "PageRepository.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	db := getDbConnection(file.Name())
	pr := NewPageRepository(db)
	pr.CreatePagesTable()

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
	updatedPage := Page{
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
	if want, got := 1, len(pages); want != got {
		t.Fatalf("Expected %d page, got %d", want, got)
	}

	updated := pages[0]
	if updated.ExtId != updatedPage.ExtId ||
		updated.Date != updatedPage.Date ||
		updated.Html != updatedPage.Html ||
		updated.Title != updatedPage.Title ||
		updated.TranslatedTitle != updatedPage.TranslatedTitle ||
		updated.Slug != updatedPage.Slug ||
		updated.Url != updatedPage.Url ||
		updated.RefUrl != updatedPage.RefUrl ||
		updated.Tags != updatedPage.Tags ||
		updated.Translated != updatedPage.Translated ||
		updated.Published != updatedPage.Published {
		t.Fatalf("Expected %v, got %v", updatedPage, updated)
	}
}

func TestPageRepositoryFindByDateReturnNothing(t *testing.T) {
	file, err := os.CreateTemp("", "PageRepository.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	db := getDbConnection(file.Name())

	pr := NewPageRepository(db)

	pr.CreatePagesTable()

	pages := pr.FindByDate("2025-01-01")

	if want, got := 0, len(pages); want != got {
		t.Fatalf("Expected %d page, got %d", want, got)
	}
}
