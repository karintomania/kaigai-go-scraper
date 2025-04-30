package scrape

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"text/template"

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

type GenerateArticle struct {
	pr        *db.PageRepository
	cr        *db.CommentRepository
	getImage  func() string
	getColour func() string
}

func NewGenerateArticle(
	pr *db.PageRepository,
	cr *db.CommentRepository,
) *GenerateArticle {
	ag := &GenerateArticle{
		pr:        pr,
		cr:        cr,
		getImage:  defaultGetImage,
		getColour: defaultGetColour,
	}

	return ag
}

func NewTestGenerateArticle(
	pr *db.PageRepository,
	cr *db.CommentRepository,
	getImage func() string,
	getColour func() string,
) *GenerateArticle {
	ag := &GenerateArticle{
		pr:        pr,
		cr:        cr,
		getImage:  getImage,
		getColour: getColour,
	}

	return ag
}

func defaultGetColour() string {
	colours := []string{"#38d3d3", "#ff5733", "#45d325", "#785bff", "#ff33a1", "#ff5c5c"}

	return colours[rand.Intn(len(colours))]
}

func (ag *GenerateArticle) run(dateStr string) error {

	pages := ag.pr.FindByDate(dateStr)

	for _, page := range pages {
		comments := ag.cr.FindByPageId(page.Id)

		article, err := ag.generateArticle(dateStr, &page, comments)
		if err != nil {
			return err
		}

		file, err := ag.getPaths(dateStr, page.Slug)
		if err != nil {
			return err
		}

		slog.Info("generating article", slog.String("path", file.Name()))

		if _, err := file.WriteString(article); err != nil {
			return err
		}
	}

	return nil
}

func (ag *GenerateArticle) getPaths(dateStr, slug string) (*os.File, error) {
	folderName := fmt.Sprintf("%s_%s",
		strings.ReplaceAll(dateStr, "-", "_"),
		slug,
	)

	outputFolder := fmt.Sprintf(
		"%s/%s/%s",
		common.GetEnv("output_article_folder"),
		dateStr,
		folderName,
	)

	if err := os.MkdirAll(outputFolder, 0774); err != nil {
		return nil, err
	}

	outputFile := fmt.Sprintf("%s/index.md", outputFolder)

	file, err := os.Create(outputFile)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (ag *GenerateArticle) generateArticle(
	dateStr string,
	page *db.Page,
	comments []db.Comment,
) (string, error) {
	commentsFiltered := make([]db.Comment, 0, 100)

	lowestScore := common.GetEnvInt("lowest_comment_score")
	minimumColourScore := common.GetEnvInt("minimum_colour_score")

	for i, _ := range comments {
		if comments[i].Score < lowestScore {
			// skip low score comments
			continue
		}

		if comments[i].Score >= minimumColourScore {
			comments[i].Colour = ag.getColour()
		}

		commentsFiltered = append(commentsFiltered, comments[i])
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

	err := frontmatterTmpl.Execute(&buf, fs)
	if err != nil {
		return "", fmt.Errorf("failed to generate frontmatter template, %w", err)
	}

	as := ArticleStruct{
		Title:           page.TranslatedTitle,
		Url:             page.Url,
		TemplateRefLink: TEMPLATE_REF_LINK,
		Comments:        commentsFiltered,
	}

	chunk := common.GetEnvInt("comment_fold_chunk_num")

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
					return (i%chunk == chunk-1 || i == len(commentsFiltered)-1) && i > chunk
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
	formatted := strings.ReplaceAll(src, "-", "/")
	formatted = strings.ReplaceAll(formatted, "T", " ")

	return formatted
}
