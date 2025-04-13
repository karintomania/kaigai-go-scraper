package app

import (
	"reflect"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func TestGetTopLinks(t *testing.T) {
	linkJsons := []JsonLink{
		{Id: "1", Points: 100, LinkText: "Link 1", Link: "http://link1.com", Dead: false},
		{Id: "2", Points: 200, LinkText: "Link 2", Link: "http://link2.com", Dead: false},
		{Id: "3", Points: 50, LinkText: "Link 3", Link: "http://link3.com", Dead: true},
		{Id: "4", Points: 150, LinkText: "Link 4", Link: "http://link4.com", Dead: false},
		{Id: "5", Points: 100, LinkText: "Link 5", Link: "http://link5.com", Dead: false},
	}

	expected := []db.Link{
		{ExtId: "2", Date: "2025-04-01", URL: "http://link2.com", Title: "Link 2", Scraped: false},
		{ExtId: "4", Date: "2025-04-01", URL: "http://link4.com", Title: "Link 4", Scraped: false},
		{ExtId: "1", Date: "2025-04-01", URL: "http://link1.com", Title: "Link 1", Scraped: false},
	}

	result := getTopLinks(linkJsons, 3, "2025-04-01")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
