package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

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

type PublishHandler struct {
	push pushFunc
	pr   *db.PageRepository
}

func NewPublishHandler(push pushFunc, pr *db.PageRepository) *PublishHandler {
	return &PublishHandler{
		push: push,
		pr:   pr,
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

	output, err := ph.push()
	if err != nil {
		result = fmt.Sprintf("Something went wrong: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		failed = true
	}

	if err := ph.updatePages(); err != nil {
		result = fmt.Sprintf("Something went wrong: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		failed = true
	}

	if !failed {
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
	slog.Info("pr", "pr", ph.pr)
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
