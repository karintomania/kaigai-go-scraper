package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/karintomania/kaigai-go-scraper/cmd/tweets"
	"github.com/karintomania/kaigai-go-scraper/db"
)

const PUBLISH_TEMPLATE = `<html>
<header>
	{{.Header}}
</header>
<body>
	<p>{{.Result}}</p>
	<a href="/">Back to home</a>
</body>
</html>
`

type pushFunc func() (string, error)
type scheduleTweetFunc func(string, []int) error

type PublishHandler struct {
	push     pushFunc
	schedule scheduleTweetFunc
	pr       *db.PageRepository
}

func NewPublishHandler(
	push pushFunc,
	pr *db.PageRepository,
	tr *db.TweetRepository,
) *PublishHandler {
	st := tweets.NewScheduleTweetsCmd(pr, tr)

	return &PublishHandler{
		push: push,
		pr:   pr,
		schedule: func(dateStr string, pageIds []int) error {
			return st.Run(dateStr, pageIds)
		},
	}
}

func (ph *PublishHandler) handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Publish handler called", "url", r.URL, "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	failed := false

	var result string

	// push git change
	output, err := ph.push()
	if err != nil {
		result = fmt.Sprintf("Something went wrong pushing git: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		failed = true
	}

	if err := ph.updatePages(); err != nil {
		result = fmt.Sprintf("Something went wrong during updating pages: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		failed = true
	}

	if err := ph.scheduleTweet(r); err != nil {
		result = fmt.Sprintf("Something went wrong during scheduling tweet: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		failed = true
	}

	if !failed {
		w.WriteHeader(http.StatusCreated)
		result = fmt.Sprintf("Success: %s", output)
	}

	tmpl := template.Must(template.New("push").Parse(PUBLISH_TEMPLATE))

	if err := tmpl.Execute(
		w,
		struct {
			Result string
			Header string
		}{
			result,
			HEADER_TEMPLATE,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *PublishHandler) updatePages() error {
	pages := ph.pr.FindUnpublished()

	for _, page := range pages {
		page.Published = true
		if err := ph.pr.Update(&page); err != nil {
			slog.Error("Error updating page", "page", page, "error", err)
			return err
		}
	}
	return nil
}

func (ph *PublishHandler) scheduleTweet(r *http.Request) error {
	defer r.Body.Close()

	dateStr := time.Now().Format("2006-01-02")

	if err := r.ParseForm(); err != nil {
		return err
	}

	pageIds := make([]int, 0, 10)

	slog.Info("page_ids", "page_ids", r.Form["page_ids"])

	for _, pageId := range r.Form["page_ids"] {
		pageIdInt, err := strconv.Atoi(pageId)
		if err != nil {
			return err
		}

		pageIds = append(pageIds, pageIdInt)
	}

	err := ph.schedule(dateStr, pageIds)

	if err != nil {
		return err
	}

	return nil
}
