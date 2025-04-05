package main

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

	db := getDbConnection(file.Name())

	lr := NewLinkRepository(db)

	lr.CreateLinksTable()

	link := Link{
		ExtId: "12345",
		Date:  "2025-01-01",
		URL:   "https://example.com/12345",
	}

	lr.Insert(link)

	links := lr.FindByDate("2025-01-01")

	if want, got := 1, len(links); want != got {
		t.Fatalf("Expected %d link, got %d", want, got)
	}

	created := links[0]

	if created.ExtId != link.ExtId ||
		created.Date != link.Date ||
		created.URL != link.URL {
		t.Fatalf("Expected %v, got %v", link, created)
	}
}

func TestLinkRepositoryFindByIdReturnNothing(t *testing.T) {
	file, err := os.CreateTemp("", "LinkRepository.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	db := getDbConnection(file.Name())

	lr := NewLinkRepository(db)

	lr.CreateLinksTable()

	links := lr.FindByDate("2025-01-01")

	if want, got := 0, len(links); want != got {
		t.Fatalf("Expected %d link, got %d", want, got)
	}
}
