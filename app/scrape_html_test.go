package app

import (
	"reflect"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
)

func TestScrapeHtml(t *testing.T) {
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
                        <span class="age" title="2025-02-07T06:23:36 1738909416">
                        <div class="commtext">Test comment 1</div>
                    </td>
                </tr>
                <tr class="athing comtr" id="456">
                    <td class="ind" indent="1"></td>
                    <td><a class="togg clicky" n="50"></a></td>
                    <td>
                        <a class="hnuser">user2</a>
                        <span class="age" title="2024-03-08T07:24:37 1738909416">
                        <div class="commtext">Test comment 2</div>
                    </td>
                </tr>
            </table>
        </body>
    </html>`

	page := &db.Page{
		Id:    99,
		ExtId: "1",
		Title: "Test Title, Which is Made !WAY! Too long intentionally to test the slug",
		Date:  "2025-01-01",
		Html:  htmlContent,
	}

	page, comments := getPageAndComments(page)

	t.Logf("Comments: %v", comments)

	if want, got := "test_title_which_is_made_way_too_long", page.Slug; want != got {
		t.Errorf("expected slug to be %s, but got %s", want, got)
	}
	
	for i, want := range []db.Comment{
		{
			ExtCommentId: "123",
			PageId:       page.Id,
			UserName:     "user1",
			Content:      "Test comment 1",
			Indent:       0,
			Reply:        100,
		},
		{
			ExtCommentId: "456",
			PageId:       page.Id,
			UserName:     "user2",
			Content:      "Test comment 2",
			Indent:       1,
			Reply:        50,
		},
	} {

		if got := comments[i]; !reflect.DeepEqual(want, got) {
			t.Errorf(
				"comment %d expected to be %v, but got %v",
				i,
				want,
				got,
			)
		}
	}

}
