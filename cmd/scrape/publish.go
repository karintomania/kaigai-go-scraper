package scrape

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

type Publish struct {
	pr            *db.PageRepository
	runGitCommand common.RunGitCommandFunc
}

func NewPublish(pr *db.PageRepository) *Publish {
	return &Publish{
		pr:            pr,
		runGitCommand: common.RunGitCommand,
	}
}

func NewPublishWithRunGitCommand(
	pr *db.PageRepository,
	runGitCommand common.RunGitCommandFunc,
) *Publish {
	return &Publish{
		pr:            pr,
		runGitCommand: runGitCommand,
	}
}

func (p *Publish) run(dateStr string) error {
	slog.Info("Publishing articles")
	srcRoot := common.GetEnv("output_article_folder")
	destRoot := common.GetEnv("article_destination")

	src := srcRoot + "/" + dateStr

	monthStr := strings.Replace(dateStr, "-", "_", 1)[0:7]
	dest := destRoot + "/" + monthStr

	// TODO: select only unpublished
	pages := p.pr.FindByDate(dateStr)

	// copy files
	for _, page := range pages {
		if err := p.copy(src, dest, dateStr, &page); err != nil {
			return err
		}

		page.Published = true
		if err := p.pr.Update(&page); err != nil {
			return err
		}

	}

	slog.Info("test")

	// run git commit
	if err := p.commit(dateStr); err != nil {
		return err
	}
	// run git push
	return nil
}

func (p *Publish) copy(src, dest, dateStr string, page *db.Page) error {
	underscoreDate := strings.ReplaceAll(dateStr, "-", "_")
	srcFilePath := fmt.Sprintf("%s/%s_%s/index.md", src, underscoreDate, page.Slug)
	destFolder := fmt.Sprintf("%s/%s_%s", dest, underscoreDate, page.Slug)
	destFilePath := fmt.Sprintf("%s/index.md", destFolder)

	slog.Info("copying file", "src", srcFilePath, "dest", destFilePath)

	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
		return err
	}

	in, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// TODO: use io.Copy if possible. Currently it shows error: copy_file_range: bad file descriptor
	content, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	if _, err := out.Write(content); err != nil {
		return err
	}

	return nil
}

func (p *Publish) commit(dateStr string) error {
	slog.Info("committing changes", "date", dateStr)

	options := []string{"add", "."}

	if _, err := p.runGitCommand(options); err != nil {
		return err
	}

	options = []string{"commit", "-m", fmt.Sprintf("add %s", dateStr)}

	if _, err := p.runGitCommand(options); err != nil {
		return err
	}

	return nil
}
