package scrape

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestScrapeHtml(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	lr := db.NewLinkRepository(dbConn)
	pr := db.NewPageRepository(dbConn)
	cr := db.NewCommentRepository(dbConn)
	scrapeHtml := NewScrapeHtml(lr, pr, cr)

	htmlContent := `<html>
        <head>
            <title>Sample Title | Example</title>
            <link rel="canonical" href="https://example.com/main-url" />
        </head>
        <body>
            <a href="https://example.com/ref-url">Sample Title</a>
            <table>
                <tr class="athing comtr" id="123">
                    <td class="ind" indent="0"></td>
                    <td><a class="togg clicky" n="100"></a></td>
                    <td>
                        <a class="hnuser">user1</a>
                        <span class="age" title="2025-01-02T03:04:05 1738909416">
                        <div class="commtext">Test comment 1</div>
                    </td>
                </tr>
                <tr class="athing comtr" id="456">
                    <td class="ind" indent="1"></td>
                    <td><a class="togg clicky" n="50"></a></td>
                    <td>
                        <a class="hnuser">user2</a>
                        <span class="age" title="2024-05-19T04:06:08 1738909416">
                        <div class="commtext">Test comment 2</div>
                    </td>
                </tr>
            </table>
        </body>
    </html>`

	date := "2025-01-01"
	page := &db.Page{
		Id:    99,
		ExtId: "1",
		Title: "Test Title, Which is Made !WAY! Too long intentionally to test the slug",
		Date:  date,
		Html:  htmlContent,
	}

	pr.Insert(page)

	t.Run("scrapePages scrape correct info", func(t *testing.T) {
		scrapeHtml.scrapePages(date)

		resultPage := pr.FindByDate(date)[0]
		resultComments := cr.FindByPageId(page.Id)

		require.Equal(t, "test_title_which_is_made_way_too_long", resultPage.Slug)

		for i, want := range []db.Comment{
			{
				ExtCommentId: "123",
				PageId:       page.Id,
				UserName:     "user1",
				Content:      "Test comment 1",
				Indent:       0,
				Reply:        100,
				CommentedAt:  "2025-01-02T03:04:05",
			},
			{
				ExtCommentId: "456",
				PageId:       page.Id,
				UserName:     "user2",
				Content:      "Test comment 2",
				Indent:       1,
				Reply:        50,
				CommentedAt:  "2024-05-19T04:06:08",
			},
		} {
			require.Equal(t, want.ExtCommentId, resultComments[i].ExtCommentId)
			require.Equal(t, want.PageId, resultComments[i].PageId)
			require.Equal(t, want.UserName, resultComments[i].UserName)
			require.Equal(t, want.Content, resultComments[i].Content)
			require.Equal(t, want.Indent, resultComments[i].Indent)
			require.Equal(t, want.Reply, resultComments[i].Reply)
			require.Equal(t, want.CommentedAt, resultComments[i].CommentedAt)
		}
	})
	t.Run("scrapePages insert only new comments", func(t *testing.T) {
		// run scrape pages again
		scrapeHtml.scrapePages(date)
		resultComments := cr.FindByPageId(page.Id)

		require.Equal(t, 2, len(resultComments))
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
