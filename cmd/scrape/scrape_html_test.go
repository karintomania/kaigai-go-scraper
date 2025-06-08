package scrape

import (
	"fmt"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

const (
	htmlWithOp = `<html>
        <head>
            <title>Sample Title | Example</title>
            <link rel="canonical" href="https://example.com/main-url" />
        </head>
        <body>
			<table class="fatitem">
				<tr><td>
				<a class="hnuser">op name</a>
				<span class="age" title="2025-01-01T02:03:04 1738909416">
				<div class="toptext">This is a comment from OP</div>
				</td></tr>
			</table>
            <a href="https://example.com/ref-url">Sample Title</a>
            <table>
                <tr class="athing comtr" id="101">
                    <td class="ind" indent="1"></td>
                    <td><a class="togg clicky" n="100"></a></td>
                    <td>
                        <a class="hnuser">user1</a>
                        <span class="age" title="2025-01-01T01:02:03 1738909416">
                        <div class="commtext">Test comment 1</div>
                    </td>
                </tr>
                <tr class="athing comtr" id="102">
                    <td class="ind" indent="2"></td>
                    <td><a class="togg clicky" n="200"></a></td>
                    <td>
                        <a class="hnuser">user2</a>
                        <span class="age" title="2025-01-02T01:02:03 1738909416">
                        <div class="commtext">Test comment 2</div>
                    </td>
                </tr>
            </table>
        </body>
    </html>`

	htmlWithoutOp = `<html>
        <head>
            <title>Sample Title | Example</title>
            <link rel="canonical" href="https://example.com/main-url" />
        </head>
        <body>
            <a href="https://example.com/ref-url">Sample Title</a>
            <table>
                <tr class="athing comtr" id="101">
                    <td class="ind" indent="1"></td>
                    <td><a class="togg clicky" n="100"></a></td>
                    <td>
                        <a class="hnuser">user1</a>
                        <span class="age" title="2025-01-01T01:02:03 1738909416">
                        <div class="commtext">Test comment 1</div>
                    </td>
                </tr>
                <tr class="athing comtr" id="102">
                    <td class="ind" indent="2"></td>
                    <td><a class="togg clicky" n="200"></a></td>
                    <td>
                        <a class="hnuser">user2</a>
                        <span class="age" title="2025-01-02T01:02:03 1738909416">
                        <div class="commtext">Test comment 2</div>
                    </td>
                </tr>
            </table>
        </body>
    </html>`
)

func TestScrapeHtml(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	common.MockEnv("max_comments_num", "10")
	common.MockEnv("max_reply_per_comment_num", "5")

	lr := db.NewLinkRepository(dbConn)
	pr := db.NewPageRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)
	scrapeHtml := NewScrapeHtml(lr, pr, cr)

	wantOpComment := db.Comment{
		ExtCommentId: "op_12",
		UserName:     "op name",
		Content:      "This is a comment from OP",
		Indent:       0,
		Reply:        9999,
		CommentedAt:  "2025-01-01T02:03:04",
	}

	wantComments := make([]db.Comment, 0, 3)

	for i := 1; i < 3; i++ {
		wantComments = append(wantComments, db.Comment{
			ExtCommentId: fmt.Sprintf("10%d", i),
			UserName:     fmt.Sprintf("user%d", i),
			Content:      fmt.Sprintf("Test comment %d", i),
			Indent:       i,
			Reply:        100 * i,
			CommentedAt:  fmt.Sprintf("2025-01-0%dT01:02:03", i),
		})
	}

	date := "2025-01-01"

	t.Run("scrapePages scrape correct info", func(t *testing.T) {
		dataTable := []struct {
			Name         string
			Html         string
			WantComments []db.Comment
		}{
			{"Article without OP", htmlWithoutOp, wantComments},
			{"Article with OP", htmlWithOp, append([]db.Comment{wantOpComment}, wantComments...)},
		}

		for _, data := range dataTable {
			pr.Truncate()
			cr.Truncate()

			page := &db.Page{
				Id:    99,
				ExtId: "12",
				Title: "Test Title, Which is Made !WAY! Too long intentionally to test the slug",
				Date:  date,
				Html:  data.Html,
			}

			pr.Insert(page)
			err := scrapeHtml.scrapePages(date)
			require.NoError(t, err)

			resultPage := pr.FindByDate(date)[0]
			resultComments := cr.FindByPageId(page.Id)

			require.Equal(t, "test_title_which_is_made_way_too_long", resultPage.Slug)

			require.Equal(t, len(data.WantComments), len(resultComments))

			for i, want := range data.WantComments {
				require.Equal(t, want.ExtCommentId, resultComments[i].ExtCommentId)
				require.Equal(t, page.Id, resultComments[i].PageId)
				require.Equal(t, want.UserName, resultComments[i].UserName)
				require.Equal(t, want.Content, resultComments[i].Content)
				require.Equal(t, want.Indent, resultComments[i].Indent)
				require.Equal(t, want.Reply, resultComments[i].Reply)
				require.Equal(t, want.CommentedAt, resultComments[i].CommentedAt)
			}

			// run scrape pages again and see comments don't duplicate
			err = scrapeHtml.scrapePages(date)
			require.NoError(t, err)

			resultCommentsOnSecondRun := cr.FindByPageId(page.Id)

			require.Equal(t, len(data.WantComments), len(resultCommentsOnSecondRun))
		}
	})

	t.Run("scrapePages insert only new comments", func(t *testing.T) {
	})
}

// Test pruning of comments
// - 1 (re:5)
//   - 2 (re:0)
//   - 3 (re:2)
//     - 4 (re:1)
//       - 5 (re:0)
//     - 6 (re:0)
// - 7 (re:1)
//   - 8 (re:0)
// - 9 (re:0)
// - 10 (re:0)

func TestSelectRelevantComments(t *testing.T) {
	maxCommentNum := 6
	maxChildCommentNum := 3

	comments := []db.Comment{
		{
			Id:     1,
			Indent: 0,
			Reply:  5,
		},
		{
			Id:     2,
			Indent: 1,
			Reply:  0,
		},
		{
			Id:     3,
			Indent: 1,
			Reply:  2,
		},
		{
			Id:     4,
			Indent: 2,
			Reply:  1,
		},
		{
			Id:     5,
			Indent: 3,
			Reply:  0,
		},
		{
			Id:     6,
			Indent: 2,
			Reply:  0,
		},
		{
			Id:     7,
			Indent: 0,
			Reply:  1,
		},
		{
			Id:     8,
			Indent: 1,
			Reply:  0,
		},
		{
			Id:     9,
			Indent: 0,
			Reply:  0,
		},
		{
			Id:     10,
			Indent: 0,
			Reply:  0,
		},
	}

	wantIds := []int{1, 2, 3, 4, 7, 8}

	result := selectRelevantComments(comments, maxCommentNum, maxChildCommentNum)

	if want, got := len(wantIds), len(result); want != got {
		t.Errorf("expected %d elements, but got %d", want, got)
	}

	for i, gotComment := range result {
		if want, got := wantIds[i], gotComment.Id; want != got {
			t.Errorf("expected %d, but got %d", want, got)
		}
	}
}
