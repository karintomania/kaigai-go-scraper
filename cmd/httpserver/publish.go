package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
)

const PUBLISH_TEMPLATE = `<html>
<body>
	{{.Result}}
</body>
</html>
`

type PublishHandler struct {
	push func() error
}

func NewPublishHandler(push func() error) *PublishHandler {
	return &PublishHandler{push: push}
}

func (ph *PublishHandler) handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Publish handler called", "url", r.URL, "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result string
	if err := ph.push(); err != nil {
		result = fmt.Sprintf("Something went wrong: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result = "Success"
	}

	tmpl := template.Must(template.New("push").Parse(PUBLISH_TEMPLATE))

	if err := tmpl.Execute(
		w,
		struct{ Result string }{result},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
