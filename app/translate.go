package app

import (
	"fmt"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/db"
)

const (
	COMMENTS_CHUNK_NUM = 10
	MAX_COMMENTS_NUM   = 90
	MAX_RETRIES        = 5
	PROMPT_SUMMARY     = `このyamlは"{title}"って記事についたコメントです
"""
comments:
{comment_str}
"""

このyamlから以下のフォーマットのJsonをこれから述べるルールにしたがって生成して。

# idフィールド
jsonのidはyamlのidをそのまま使ってください

# commentフィールド
jsonのcommentの値はyamlのcommenをルールに従って翻訳・要約して。
- カジュアルなタメ口の日本語に翻訳して
- 翻訳が200字以上になるコメントは200字以内に要約して
- 記号は使わず翻訳に含めるときは全角の記号に変換して。**特に'と", \はJsonが壊れるので’と”に変換するか省略して**
- 固有名詞は英語のままにして
- >で始まる引用は、先頭に全角の＞をつけ翻訳した引用部を全角の”でくくることで地の文と違いがわかるようにして
- 改行は\\nではなく<br>にして

# scoreフィールド
jsonのscoreはそれぞれのコメントを以下のルールに従って採点して。
- 点数は0から100
- コメントが役に立つ、面白いなら高得点、情報量が少ないものや関係のないコメントは低い点にする

# Jsonに変換
- 返事はバックティックで囲わず平文のJsonで返事して
- id, comment, scoreは必須項目。
- 入力と出力のコメント数が同じになるようにして
- フィールド内の"は全角に変換して
- Jsonのバリデーションをしてから返事して
- 以下のJSONのとおりに出力して
"""
{"comments": [{"id": "1234", "comment": "翻訳コメント", "score": 90}]}
"""`

	PROMPT_TITLE = `次のHacker Newsの記事タイトル「%s」について以下のタスクをしてください。
ステップ１：タイトルを以下のルールにしたがって日本語に訳してください。
翻訳ルール
- まとめサイトの記事っぽくして
- 興味を引くタイトルにして
- 原文にない情報を付け加えないで
- タイトル中にShow HN, Ask HNがあるとき、それは訳に含めないで
- ！や？は必要な際には使っていいけど、「」【】, "", ''などの記号は使わないで。
- 以下の記事に対するコメントを文脈を捉えるのに使って。
"""
%s
"""
ステップ２：タイトルから関連タグを最大５個考えて

結果は以下のJsonのとおりにして、バックティックなどを使わずPlain textで答えてください。
"""
{
  "title": "翻訳後のタイトル",
  "tags": ["HTML", "AI", "プログラミング"]
}
"""`
)

func translate(
	dateString string,
	pageRepository *db.PageRepository,
	commentRepository *db.CommentRepository,
) error {
	pages := pageRepository.FindByDate(dateString)

	for _, page := range pages {
		comments := commentRepository.FindByPageId(page.Id)

		translatePageAndComments(&page, comments)

		// update page, comments
		pageRepository.Update(&page)

		for _, comment := range comments {
			commentRepository.Update(&comment)
		}
	}

	return nil
}

func translatePageAndComments(page *db.Page, comments []db.Comment) {
	// translate page

	// translate comments
}

func translatePage(page *db.Page, comments []db.Comment) {
	// build prompt
	var sb strings.Builder
	for _, comment := range comments {
		fmt.Fprintln(&sb, comment.Content)
	}

	commentsStr := sb.String()

	prompt := fmt.Sprintf(PROMPT_TITLE, page.Title, commentsStr)

	fmt.Println(prompt)
}
