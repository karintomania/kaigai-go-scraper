package db

import (
	"os"
	"testing"
)

func TestLinkRepositoryInsert(t *testing.T) {
	file, err := os.CreateTemp("", "LinkRepository.sql")

	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(file.Name())

	db := GetDbConnection(file.Name())
	defer db.Close()

	lr := NewLinkRepository(db)

	lr.CreateLinksTable()

	link := Link{
		ExtId:   "12345",
		Date:    "2025-01-01",
		URL:     "https://example.com/12345",
		Title:   "Example Title",
		Scraped: false,
	}

	lr.Insert(&link)

	links := lr.FindByDate("2025-01-01")

	if want, got := 1, len(links); want != got {
		t.Fatalf("Expected %d link, got %d", want, got)
	}

	created := links[0]

	if created.ExtId != link.ExtId ||
		created.Date != link.Date ||
		created.URL != link.URL ||
		created.Title != link.Title ||
		created.Scraped != link.Scraped {
		t.Fatalf("Expected %v, got %v", link, created)
	}
}

func TestLinkRepositoryFindByIdReturnNothing(t *testing.T) {
	file, err := os.CreateTemp("", "LinkRepository.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	db := GetDbConnection(file.Name())
	defer db.Close()

	lr := NewLinkRepository(db)

	lr.CreateLinksTable()

	links := lr.FindByDate("2025-01-01")

	if want, got := 0, len(links); want != got {
		t.Fatalf("Expected %d link, got %d", want, got)
	}
}
