package app

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/karintomania/kaigai-go-scraper/common"
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
{{ "{{" }}<matomeQuote body="{{ $comment.TranslatedContent }}" userName="{{ $comment.UserName }}" createdAt="{{ (formatCommentedAt $comment.CommentedAt) }}" color="{{ $comment.Colour }}">{{ "}}" }}
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
	getColour func() string
}

func NewArticleGenerator(
	pr *db.PageRepository,
	cr *db.CommentRepository,
) *ArticleGenerator {
	ag := &ArticleGenerator{
		pr:       pr,
		cr:       cr,
		getImage: defaultGetImage,
		getColour: defaultGetColour,
	}

	return ag
}

func NewTestArticleGenerator(
	pr *db.PageRepository,
	cr *db.CommentRepository,
	getImage func() string,
	getColour func() string,
) *ArticleGenerator {
	ag := &ArticleGenerator{
		pr:       pr,
		cr:       cr,
		getImage: getImage,
		getColour: getColour,
	}

	return ag
}

func defaultGetColour() string {
	colours := []string{"#38d3d3", "#ff5733", "#45d325", "#785bff", "#ff33a1", "#ff5c5c"}
	return colours[rand.Intn(len(colours))]
}

func (ag *ArticleGenerator) generateArticles(dateStr string) error {

	pages := ag.pr.FindByDate(dateStr)

	for _, page := range pages {
		comments := ag.cr.FindByPageId(page.Id)

		article, err := ag.generateArticle(dateStr, &page, comments)

		if err != nil { return err }

		outputFolder := fmt.Sprintf(
			"%s/%s/%s",
			common.GetEnv("output_article_folder"),
			dateStr,
			page.Slug,
			)
		if err := os.MkdirAll(outputFolder, 0774); err != nil {
			return err
		}

		outputFile := fmt.Sprintf("%s/index.md", outputFolder)

		file, err := os.Create(outputFile)
		if err != nil { return err }


		if _, err := file.WriteString(article); err != nil {
			return err
		}
	}

	return nil
}

func (ag *ArticleGenerator) generateArticle(
	dateStr string,
	page *db.Page,
	comments []db.Comment,
) (string, error) {
	minimumColourScore, err := strconv.Atoi(common.GetEnv("minimum_colour_score"))
	if err != nil {
		return "", err
	}

	for i, _ := range comments {
		if comments[i].Score >=  minimumColourScore{
			comments[i].Colour = ag.getColour()
		}
	}

	fs := FrontmatterStruct{
		Date:  dateStr + "T00:00:00",
		Month: strings.Replace(dateStr[0:7], "-", "/", -1),
		Title: page.TranslatedTitle,
		Tags:  `"` + strings.Replace(page.Tags, `,`, `", "`, -1) + `"`,
		Image: ag.getImage(),
	}

	frontmatterTmpl := template.Must(template.New("frontmatter").Parse(TEMPLATE_FRONTMATTER))
	var buf bytes.Buffer

	err = frontmatterTmpl.Execute(&buf, fs)
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
				"formatCommentedAt": func(str string) string {
					return formatCommentedAt(str)
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

func formatCommentedAt(src string) string {
	date, err := time.Parse(db.Rfc3339Milli, src)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	formatted := date.Format("2006/01/02 15:04:05")

	return formatted
}
