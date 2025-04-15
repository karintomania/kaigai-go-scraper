package external

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/db"
)


func CallHackerNews(link *db.Link) (string, io.ReadCloser) {
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

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("api call failed. status code: %d, URL: %s\n", resp.StatusCode, url)
	}

	return url, resp.Body
}

func CallHckrNews(date string) io.ReadCloser {
	date = strings.Replace(date, "-", "", -1)
	url := fmt.Sprintf("https://hckrnews.com/data/%s.js", date)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("api call failed. status code: %d\n", resp.StatusCode)
	}

	return resp.Body
}
