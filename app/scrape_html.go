package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/karintomania/kaigai-go-scraper/db"
)

func scrapeHtml(
	dateString string,
	linkRepository *db.LinkRepository,
	pageRepository *db.PageRepository,
	commentRepository *db.CommentRepository,
) error {
	links := linkRepository.FindByDate(dateString)

	for _, link := range links {

		if err := downloadHtml(&link, dateString, pageRepository); err != nil {
			return err
		}

		// mark link as scraped
		link.Scraped = true
		linkRepository.Update(&link)

	}

	pages := pageRepository.FindByDate(dateString)

	for _, page := range pages {
		_, comments := getPageAndComments(&page)
		pageRepository.Update(&page)

		for _, comment := range comments {
			commentRepository.Insert(&comment)
		}
	}

	return nil
}

func downloadHtml(link *db.Link, dateString string, pageRepository *db.PageRepository) error {
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

// scrape info from the HTML
// update page and return comments
func getPageAndComments(page *db.Page) (*db.Page, []db.Comment) {
	fmt.Println("test")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Html))

	if err != nil {
		log.Fatalln(err)
	}

	comments := make([]db.Comment, 0)

	doc.Find("tr.athing.comtr").Each(func (i int, s *goquery.Selection) {

		replyStr := s.Find(".clicky").AttrOr("nn", "22")

		reply, err := strconv.Atoi(replyStr)
		fmt.Printf("reply: %v, replyStr %s\n", reply, replyStr)
		if err != nil {
			log.Fatalln(err)
		}

		indent, err := strconv.Atoi(s.Find("td.ind").AttrOr("indent", "0"))
		if err != nil {
			log.Fatalln(err)
		}

		comment := db.Comment{
			ExtCommentId: s.AttrOr("id", ""),
			PageId:       page.Id,
			UserName:     s.Find("a.hnuser").Text(),
			Content:      strings.TrimSpace(s.Find(".commtext").Text()),
			Indent:       indent,
			Reply:        reply,
		}

		fmt.Printf("comment: %v\n", comment)
		comments  = append(comments, comment)
		
	})

	return page, comments
}
