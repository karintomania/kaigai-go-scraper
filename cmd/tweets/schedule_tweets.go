package tweets

import (
	"bytes"
	"log/slog"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/karintomania/kaigai-go-scraper/db"
)

const (
	TWEET_TEMPLATE = `「{{.Title}}」に対する海外の反応をまとめました。
#海外の反応 #テックニュース

{{.PostUrl}}/{{.YearMonth}}/{{.Slug}}/`
	date_str_format = "2006-01-02"
)

type ScheduleTweetsCmd struct {
	pr *db.PageRepository
	tr *db.TweetRepository
}

func (cmd *ScheduleTweetsCmd) Run(dateStr string) error {
	pages := cmd.pr.FindByDate(dateStr)

	for _, page := range pages {
		content := cmd.createTweetContent(&page)
		scheduledDate := cmd.generateScheduleDate(page.Date)

		tweet := db.Tweet{
			PageId:      page.Id,
			Date:        page.Date,
			Content:     content,
			ScheduledAt: scheduledDate,
			Published:   false,
		}

		cmd.tr.Insert(&tweet)
	}

	return nil
}

func (cmd *ScheduleTweetsCmd) createTweetContent(page *db.Page) string {
	tmpl := template.Must(template.New("tweet").Parse(TWEET_TEMPLATE))

	var buf bytes.Buffer

	yearMonth := strings.ReplaceAll(page.Date[:7], "-", "_")

	err := tmpl.Execute(&buf, struct {
		Title     string
		YearMonth string
		Slug      string
		PostUrl   string
	}{
		page.Title,
		yearMonth,
		page.Slug,
		"https://www.kaigai-tech-matome.com/posts",
	})

	if err != nil {
		slog.Error("Failed to execute template", "error", err)
		panic(err)
	}

	return buf.String()
}

func (cmd *ScheduleTweetsCmd) generateScheduleDate(startDateStr string) string {
	date, err := time.Parse(date_str_format, startDateStr)
	if err != nil {
		panic(err)
	}

	rnd := rand.Intn(7) + 1

	date = date.AddDate(0, 0, rnd)

	return date.Format(date_str_format)
}
