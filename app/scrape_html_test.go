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
                    <a class="togg clicky" n="100"></a>
                    <td>
                        <a class="hnuser">user1</a>
                        <span class="age" title="2025-02-07T06:23:36 1738909416">
                        <div class="commtext">Test comment 1</div>
                    </td>
                </tr>
                <tr class="athing comtr" id="456">
                    <td class="ind" indent="1"></td>
                    <a class="togg clicky" n="50"></a>
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
		Id:    1,
		ExtId: "1",
		Date:  "2025-01-01",
		Html:  htmlContent,
	}

	page, comments := getPageAndComments(page)

	t.Logf("Comments: %v", comments)
	
	for i, want := range []db.Comment{
		{
			Id:           1,
			ExtCommentId: "123",
			PageId:       page.Id,
			UserName:     "user1",
			Content:      "test comment 1",
			Indent:       0,
			Reply:        100,
		},
		{
			Id:           2,
			ExtCommentId: "456",
			PageId:       page.Id,
			UserName:     "user2",
			Content:      "test comment 2",
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
