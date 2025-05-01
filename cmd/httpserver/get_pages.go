package httpserver

import (
	"net/http"
	"text/template"

	"github.com/karintomania/kaigai-go-scraper/db"
)

const GET_PAGE_TEMPLATE = `<html>
<body>
{{range $key, $value := .DatePagesMap}}
<h1>{{$key}}</h1>
{{range $value}}
	<h2>{{.TranslatedTitle}}</h2>
{{end}}
{{end}}
	<form method="POST" action="/publish">
		<input type="submit" value="Publish" />
	</form>
</body>
</html>
`

type DatePagesMap map[string][]db.Page

type GetPageHandler struct {
	pr      *db.PageRepository
	dateStr string
}

func NewGetPageHandler(pr *db.PageRepository, dateStr string) *GetPageHandler {
	return &GetPageHandler{pr: pr}
}

func (gph *GetPageHandler) getPages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	datePagesMap := make(DatePagesMap)

	pages := gph.pr.FindUnpublished()

	for _, page := range pages {
		date := page.Date
		if _, ok := datePagesMap[date]; ok {
			datePagesMap[date] = append(datePagesMap[date], page)
		} else {
			datePagesMap[date] = []db.Page{page}
		}
	}

	tmpl := template.Must(template.New("get_pages").Parse(GET_PAGE_TEMPLATE))

	if err := tmpl.Execute(
		w,
		struct{ DatePagesMap DatePagesMap }{datePagesMap},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
