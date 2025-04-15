package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/db"
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

const TOP_LINK_NUM = 10

func StoreLinks(dateString string, linkRepository *db.LinkRepository) error {
	jsonLinks, err := callHckrNewsApi(dateString)
	if err != nil {
		return err
	}

	links := getTopLinks(jsonLinks, TOP_LINK_NUM, dateString)

	for _, link := range links {
		linkRepository.Insert(&link)
		fmt.Printf("Inserted link: %s\n", link.URL)
	}

	return nil
}

func callHckrNewsApi(date string) ([]JsonLink, error) {
	date = strings.Replace(date, "-", "", -1)
	url := fmt.Sprintf("https://hckrnews.com/data/%s.js", date)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api call failed. status code: %d", resp.StatusCode)
	}

	var links []JsonLink
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return nil, err
	}

	return links, nil
}

// get top n links by points. Convert LinkJson to Link
func getTopLinks(linkJsons []JsonLink, n int, dateString string) []db.Link {
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
