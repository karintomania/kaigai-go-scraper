package app

import (
	"reflect"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func TestScrapeHtml(t *testing.T) {

	t.Log("test scrape html")

	page := &db.Page{
		Id:    1,
		ExtId: "1",
		Date:  "2025-01-01",
		Html:  "<html></html>",
	}

	page, comments := getPageAndComments(page)

	if want, got := "test title", page.Title; want != got {
		t.Errorf("expected title %s, but got %s", want, got)
	}

	for i, want := range []db.Comment{
		{
			Id:           1,
			ExtCommentId: "1",
			PageId:       1,
			UserName:     "User 1",
			Content:      "test comment 1",
			Indent:       "0",
			Reply:        "1",
		},
		{
			Id:           1,
			ExtCommentId: "2",
			PageId:       2,
			UserName:     "User 2",
			Content:      "test comment 2",
			Indent:       "1",
			Reply:        "2",
		},
	} {

		if got := comments[i]; !reflect.DeepEqual(want, got) {
			t.Errorf(
				"comment %d expected to be %v, but got %v",
				i,
				want,
				got,
			)
		}
	}

}
