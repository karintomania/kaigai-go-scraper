package app

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

const MAX_COMMENTS_NUM = 100
const MAX_REPLY_PER_COMMENT_NUM = 30

func scrapeHtml(
	dateString string,
	linkRepository *db.LinkRepository,
	pageRepository *db.PageRepository,
	commentRepository *db.CommentRepository,
) error {

	// if err := downloadHtmlsAsync(dateString, linkRepository, pageRepository); err != nil {
	// 	return err
	// }

	pages := pageRepository.FindByDate(dateString)

	for _, page := range pages {
		_, comments := getPageAndComments(&page)
		pageRepository.Update(&page)

		selectedComments := selectRelevantComments(comments, MAX_COMMENTS_NUM, MAX_REPLY_PER_COMMENT_NUM)

		for _, comment := range selectedComments {
			commentRepository.Insert(&comment)
		}
	}

	return nil
}

func downloadHtmlsAsync(
	dateString string,
	linkRepository *db.LinkRepository,
	pageRepository *db.PageRepository,
) error {
	links := linkRepository.FindByDate(dateString)

	var wg sync.WaitGroup
	var errEncountered error
	var errMutex sync.Mutex

	for _, link := range links {
		wg.Add(1)
		time.Sleep(1 * time.Second)

		go func(link *db.Link) {
			err := downloadHtml(link, dateString, pageRepository)
			if err != nil {
				slog.Error(
					"Error happend while downloading HTML",
					slog.String("url", link.URL),
					slog.Any("err", err),
				)

				errMutex.Lock()
				errEncountered = err
				errMutex.Unlock()
				return
			}

			// mark link as scraped
			link.Scraped = true
			linkRepository.Update(link)

			slog.Info("HTML downloaded", slog.String("title", link.Title), slog.Int("link id", link.Id))

			wg.Done()
		}(&link)
	}
	wg.Wait()

	if errEncountered != nil {
		return fmt.Errorf("error downloading HTML, %w", errEncountered)
	}

	return nil
}

func downloadHtml(link *db.Link, dateString string, pageRepository *db.PageRepository) error {
	url, body := external.CallHackerNews(link)

	htmlBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("error on reading response from %s: %w", url, err)
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
	doc.Find("tr.athing.comtr").Each(func(i int, s *goquery.Selection) {
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

		comments = append(comments, comment)
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

// select comments by relevance
// HN already sort the comments by relevance,
// but sometimes, child comments are not really related to the main thread.
// So, this method only gets the first n-th child comments by reply number
func selectRelevantComments(
	comments []db.Comment,
	maxCommentNum int,
	maxChildCommentNum int,
) []db.Comment {
	if len(comments) < maxCommentNum {
		return comments
	}

	result := make([]db.Comment, 0)

	children := make([]db.Comment, 0)

	for i, c := range comments {
		if c.Indent == 0 {
			result = append(result, c)
		} else {
			children = append(children, c)

			// if this is the last child comment
			next := i + 1

			if next == len(comments) || comments[next].Indent == 0 {
				// prune child comments if child comments are too many
				if len(children) >= maxChildCommentNum {
					pruneChildrenComments(&children, maxChildCommentNum)
				}

				// add children to the result
				result = append(result, children...)
				// reset children
				children = make([]db.Comment, 0)
			}
		}
	}

	return result[:maxCommentNum]
}

func pruneChildrenComments(commentsPtr *[]db.Comment, maxChildCommentNum int) {
	comments := *commentsPtr

	// sort by reply number
	sort.Slice(comments, func(i, j int) bool {
		if comments[j].Reply != comments[i].Reply {
			return comments[i].Reply > comments[j].Reply
		}

		// if the reply is the same, sort by id ascend
		return comments[i].Id < comments[j].Id
	})

	// get first n-th (n = max child comment limit)
	comments = comments[:maxChildCommentNum]

	// sort by id ascend again
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].Id < comments[j].Id
	})

	// update comments
	*commentsPtr = comments
}
