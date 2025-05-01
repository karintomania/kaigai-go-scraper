package scrape

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	common.SetLogger()
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	dateStr := "2025-01-01"

	// create test files
	tmpDir, err := os.MkdirTemp(os.TempDir(), "TestPublish")
	require.NoError(t, err)

	srcRootDir := tmpDir + "/src"
	destRootDir := tmpDir + "/dest"

	common.MockEnv("output_article_folder", srcRootDir)
	common.MockEnv("article_destination", destRootDir)

	srcPageDir := fmt.Sprintf("%s/%s/%s", srcRootDir, dateStr, "/2025_01_01_test_slug")

	for _, dir := range []string{srcRootDir, destRootDir, srcPageDir} {
		err := os.MkdirAll(dir, os.ModePerm)
		require.NoError(t, err)
	}

	srcFile, err := os.Create(srcPageDir + "/index.md")
	require.NoError(t, err)

	_, err = srcFile.WriteString("test content")
	require.NoError(t, err)

	page := db.Page{
		Id:   1,
		Slug: "test_slug",
		Date: dateStr,
	}

	gitCounter := 1
	mockRunGitCommand := func(args []string) error {
		if gitCounter == 1 {
			require.Equal(t, args[0], "add")
			require.Equal(t, args[1], ".")
			gitCounter++
		} else if gitCounter == 2 {
			require.Equal(t, args[0], "commit")
			require.Equal(t, args[1], "-m")
			require.Equal(t, args[2], "add 2025-01-01")
			gitCounter++
		} else {
			t.Error("This should be called only twice")
		}
		return nil
	}

	pr := db.NewPageRepository(dbConn)

	p := NewPublishWithRunGitCommand(pr, mockRunGitCommand)

	pr.Insert(&page)

	err = p.run(dateStr)

	require.NoError(t, err)

	copied, err := os.Open(destRootDir + "/2025_01/2025_01_01_test_slug/index.md")
	require.NoError(t, err)
	copiedContent, err := io.ReadAll(copied)
	require.NoError(t, err)
	require.Equal(t, string(copiedContent), "test content")
}
