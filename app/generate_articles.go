package app

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/karintomania/kaigai-go-scraper/db"
)

const (

	// use {{ "}}" }} to escape curly braces
	TEMPLATE_PAGE = `
> {{.Title}}

引用元：[{{.Url}}]({{.Url}})

{{ range $i, $comment := .Comments }}
{{/* Show fold every 10 comments */}}
{{ if (toStartDetails $i) }}
{{ "{{" }}< details summary="もっとコメントを表示（{{getDetailPageCount $i}}）">{{ "}}" }}
{{ end }}
{{ "{{" }}<matomeQuote body="{{ $comment.TranslatedContent }}" userName="{{ $comment.UserName }}" createdAt="2025/02/02" color="{{ $comment.Colour }}">{{ "}}" }}
{{ if (toCloseDetails $i ) }}
{{ "{{" }}</details>{{ "}}" }}
{{ end }}
{{ end }}

[記事一覧へ]({{ .TemplateRefLink}})
`
	TEMPLATE_FRONTMATTER = `+++
date = '{{.Date}}'
months = '{{.Month}}'
draft = false
title = '{{.Title}}'
tags = [{{.Tags}}]
featureimage = '{{.Image}}'
+++
`
	// This string is extremely hard to escape inside template, so not included.
	TEMPLATE_REF_LINK = `{{% ref "/posts/" %}}`
)

type ArticleGenerator struct {
	pr       *db.PageRepository
	cr       *db.CommentRepository
	getImage func() string
}

func NewArticleGenerator(
	pr *db.PageRepository,
	cr *db.CommentRepository,
	options ...func(*ArticleGenerator),
) *ArticleGenerator {
	ag := &ArticleGenerator{
		pr:       pr,
		cr:       cr,
		getImage: defaultGetImage,
	}

	for _, options := range options {
		options(ag)
	}

	return ag
}

func WithGetImage(getImage func() string) func(*ArticleGenerator) {
	return func(ag *ArticleGenerator) {
		ag.getImage = getImage
	}
}

func defaultGetColour() string {
	return "#000000"
}

func (ag *ArticleGenerator) generateArticles(dateStr string) error {

	return nil
}

func (ag *ArticleGenerator) generateArticle(
	dateStr string,
	page *db.Page,
	comments []db.Comment,
) (string, error) {
	fs := FrontmatterStruct{
		Date:  dateStr + "T00:00:00",
		Month: strings.Replace(dateStr[0:7], "-", "/", -1),
		Title: page.TranslatedTitle,
		Tags:  `"` + strings.Replace(page.Tags, `,`, `", "`, -1) + `"`,
		Image: ag.getImage(),
	}

	frontmatterTmpl := template.Must(template.New("frontmatter").Parse(TEMPLATE_FRONTMATTER))
	var buf bytes.Buffer

	err := frontmatterTmpl.Execute(&buf, fs)
	if err != nil {
		return "", fmt.Errorf("failed to generate frontmatter template, %w", err)
	}

	as := ArticleStruct{
		Title:           page.TranslatedTitle,
		Url:             page.Url,
		TemplateRefLink: TEMPLATE_REF_LINK,
		Comments:        comments,
	}

	chunk := 10
	pageTmpl := template.Must(
		template.New("page").Funcs(
			template.FuncMap{
				"getDetailPageCount": func(i int) int {
					return i / chunk
				},
				"toStartDetails": func(i int) bool {
					return i%chunk == 0 && i != 0
				},
				"toCloseDetails": func(i int) bool {
					// if i > chunk, show close detail tag for
					// the last comment of the chunk or the last comment
					return (i%chunk == chunk-1 || i == len(comments)-1) && i > chunk
				},
			},
		).Parse(TEMPLATE_PAGE))

	var pageBuf bytes.Buffer
	err = pageTmpl.Execute(&pageBuf, as)
	if err != nil {
		return "", fmt.Errorf("Failed to generate page template, %w", err)
	}

	return buf.String() + pageBuf.String(), nil
}

type FrontmatterStruct struct {
	Date  string
	Month string
	Title string
	Tags  string
	Image string
}

type ArticleStruct struct {
	Title           string
	Url             string
	TemplateRefLink string
	Comments        []db.Comment
}
