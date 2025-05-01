package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
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
}

func NewPublishHandler(push pushFunc) *PublishHandler {
	return &PublishHandler{push: push}
}

func (ph *PublishHandler) handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("Publish handler called", "url", r.URL, "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var result string
	if output, err := ph.push(); err != nil {
		result = fmt.Sprintf("Something went wrong: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result = fmt.Sprintf("Success: %s", output)
	}

	tmpl := template.Must(template.New("push").Parse(PUBLISH_TEMPLATE))

	if err := tmpl.Execute(
		w,
		struct{
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
