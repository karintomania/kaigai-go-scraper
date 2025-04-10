package app

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func scrapeHtml(
	dateString string,
	linkRepository *db.LinkRepository,
	pageRepository *db.PageRepository,
) error {
	links := linkRepository.FindByDate(dateString)

	for _, link := range links {
		if err := storeHtml(link, dateString, pageRepository); err != nil {
			return err
		}
		link.Scraped = true
		linkRepository.Update(&link)
	}

	return nil
}

func storeHtml(link db.Link, dateString string, pageRepository *db.PageRepository) error {
	url := "https://news.ycombinator.com/item?id=" + link.ExtId

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error on scraping %s. Http Status %s", url, resp.Status)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error on reading response from %s, %w", url, err)
	}

	html := string(htmlBytes)

	page := db.Page{
		ExtId:  link.ExtId,
		Date:   dateString,
		Title:  link.Title,
		Html:   html,
		Url:    url,
		RefUrl: link.URL,
	}

	pageRepository.Insert(&page)

	return nil
}
