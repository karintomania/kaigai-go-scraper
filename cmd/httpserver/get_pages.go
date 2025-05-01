package httpserver

import (
	"net/http"
	"text/template"

	"github.com/karintomania/kaigai-go-scraper/db"
)



const (
	GET_PAGE_TEMPLATE = `<html>
<header>
{{.Header}}
</header>
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
		struct{
			DatePagesMap DatePagesMap
			Header string
		}{
			datePagesMap,
			HEADER_TEMPLATE,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
