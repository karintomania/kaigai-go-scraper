package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
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
	// set slug for the page
	page.Slug = getSlug(page.Title)

	// create parser
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Html))

	if err != nil {
		log.Fatalln(err)
	}

	comments := make([]db.Comment, 0)

	// loop all comments
	doc.Find("tr.athing.comtr").Each(func (i int, s *goquery.Selection) {
		// get reply
		reply, err := strconv.Atoi(s.Find("a.clicky.togg").AttrOr("n", "0"))
		if err != nil {
			log.Fatalln(err)
		}

		// get indent
		indent, err := strconv.Atoi(s.Find("td.ind").AttrOr("indent", "0"))
		if err != nil {
			log.Fatalln(err)
		}

		// create comment struct
		comment := db.Comment{
			ExtCommentId: s.AttrOr("id", ""),
			PageId:       page.Id,
			UserName:     s.Find("a.hnuser").Text(),
			Content:      strings.TrimSpace(s.Find(".commtext").Text()),
			Indent:       indent,
			Reply:        reply,
		}

		comments  = append(comments, comment)
	})

	return page, comments
}

func getSlug(str string) string {
	lower := strings.ToLower(str)

	re := regexp.MustCompile("[a-z0-9]+")

	maxWordCount := 8

	// get min(all, maxWordCount) words for slug
	words := re.FindAllString(lower, maxWordCount)

	return strings.Join(words, "_")
	
}
