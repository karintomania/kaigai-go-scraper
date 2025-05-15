package httpserver

import (
	"net/http"
	"text/template"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

const (
	GET_PAGE_TEMPLATE = `<html>
<header>
{{.Header}}
</header>
<body>
<h2><a href="{{.Host}}">Open Blog</a></h2>
{{range $key, $value := .DatePagesMap}}
<form method="POST" action="/publish">
<h1>{{$key}}</h1>
{{range $value}}
	<div>
	<h2>{{.TranslatedTitle}}</h2>
	<a href={{.RefUrl}}>{{.Title}}</a><br>
	Tweet: <input type="checkbox" name="page_ids" value="{{.Id}}" />
	</div>
{{end}}
{{else}}
	<h1>All clear!</h1>
	<p>Nothing to publish.</p>
{{end}}
<input type="submit" value="Publish" />
</form>
</body>
</html>
`
	HEADER_TEMPLATE = `<style>
body {
    font-family: sans-serif;
    margin: 20px;
    background-color: #f4f4f4;
    color: #333;
}
h1 {
    color: #0056b3;
    border-bottom: 2px solid #0056b3;
    padding-bottom: 5px;
    margin-top: 20px;
}
h2 {
    color: #333;
    margin-left: 20px;
}
form {
    margin-top: 30px;
}
input[type="submit"] {
    background-color: #007bff;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
}
input[type="submit"]:hover {
    background-color: #0056b3;
}
p {
    line-height: 1.6;
    margin-bottom: 15px;
}
a {
    color: #007bff;
    text-decoration: none;
}
a:hover {
    text-decoration: underline;
}
</style>`
)

type DatePagesMap map[string][]db.Page

type GetPageHandler struct {
	pr *db.PageRepository
}

func NewGetPageHandler(pr *db.PageRepository) *GetPageHandler {
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

		_, ok := datePagesMap[date]
		if ok {
			datePagesMap[date] = append(datePagesMap[date], page)
		} else {
			datePagesMap[date] = []db.Page{page}
		}
	}

	tmpl := template.Must(template.New("get_pages").Parse(GET_PAGE_TEMPLATE))

	if err := tmpl.Execute(
		w,
		struct {
			DatePagesMap DatePagesMap
			Header       string
			Host         string
		}{
			datePagesMap,
			HEADER_TEMPLATE,
			common.GetEnv("server_host"),
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
