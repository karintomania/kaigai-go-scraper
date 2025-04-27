package app

import (
	"encoding/json"
	"io"
	"log/slog"
	"sort"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

type JsonLink struct {
	Id       string `json:"id"`
	Points   int    `json:"points"`
	LinkText string `json:"link_text"`
	Link     string `json:"link"`
	Dead     bool   `json:"dead"`
	Source   string `json:"source"`
	// Type      string `json:"type"`
	// Submitter string `json:"submitter"`
	// Time      string `json:"time"`
	// Date      int    `json:"date"`
	// Comments  int    `json:"comments"`
}

type StoreLinks struct {
	lr             *db.LinkRepository
	callHackerNews func(string) io.ReadCloser
}

func NewStoreLinks(lr *db.LinkRepository) *StoreLinks {
	return &StoreLinks{
		lr:             lr,
		callHackerNews: external.CallHckrNews,
	}
}

func NewTestStoreLinks(lr *db.LinkRepository, callHckrNews func(string) io.ReadCloser) *StoreLinks {
	return &StoreLinks{
		lr:             lr,
		callHackerNews: callHckrNews,
	}
}

func (sl *StoreLinks) run(dateString string) error {
	slog.Info("Start storing links")

	topLinkNum := common.GetEnvInt("top_link_num")

	jsonLinks, err := sl.callHckrNewsApi(dateString)
	if err != nil {
		return err
	}

	links := sl.getTopLinks(jsonLinks, topLinkNum, dateString)

	sl.storeTopLinks(links)

	return nil
}

func (sl *StoreLinks) callHckrNewsApi(date string) ([]JsonLink, error) {
	body := sl.callHackerNews(date)

	defer body.Close()

	var links []JsonLink
	if err := json.NewDecoder(body).Decode(&links); err != nil {
		return nil, err
	}

	return links, nil
}

// get top n links by points. Convert LinkJson to Link
func (sl *StoreLinks) getTopLinks(linkJsons []JsonLink, n int, dateString string) []db.Link {
	sort.Slice(linkJsons, func(i, j int) bool {
		return linkJsons[i].Points > linkJsons[j].Points
	})

	links := make([]db.Link, 0, n)
	for i := 0; len(links) < n && i < len(linkJsons); i++ {
		linkJson := linkJsons[i]

		if linkJson.Dead {
			continue // skip dead links
		}

		link := db.Link{
			ExtId:   linkJson.Id,
			Date:    dateString,
			URL:     linkJson.Link,
			Title:   linkJson.LinkText,
			Scraped: false,
		}
		links = append(links, link)
	}

	return links
}

func (sl *StoreLinks) storeTopLinks(links []db.Link) {
	for _, link := range links {
		if !sl.lr.DoesExternalIdExist(link.ExtId) {
			sl.lr.Insert(&link)
			slog.Info("Inserted", "link", link.URL, "title", link.Title)
		}
	}
}
