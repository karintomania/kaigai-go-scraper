package scrape

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

const (
	// Move these consts to .env
	COMMENTS_CONTEXT_NUM = 10
	COMMENTS_CHUNK_NUM   = 10
	MAX_RETRIES          = 5
	PROMPT_COMMENT       = `このjsonは"%s"って記事についたコメントです
"""
%s
"""

# ゴール
このjsonのcontentを翻訳しscoreの値を計算して

# **最重要**
* JSON形式が正しいようにして
* コメントは必ず200字以内に要約してね
* 翻訳後のcommentには日本語、英語の固有名詞だけ(他の言語は使わない)があるようにして

# タスク
1.  contentの翻訳と要約:
- 各contentフィールドを200字以内のタメ口の日本語に要約して
- 記号は基本的に削除、必要なら全角に変換し、'は’、"は”、\は＼にして
- 改行は<br>を使用してね
- 対応する日本語がない英語の固有名詞（例: JavaScript, LLM）はそのまま使用しても大丈夫。
- URLは省略しないこと。URLを入れるなら200字を超えてもいい。
2. scoreの採点:
- 各コメントのscoreを0から100の範囲で採点して。
- 役に立つとか面白いコメントは高得点、情報が少ないか無関係なコメントは低得点にして。
- 入力のscoreの値は無視して。

# 出力形式
入力と同じJSON形式で出力して。`

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

type TranslatePage struct {
	pr     *db.PageRepository
	cr     *db.CommentRepository
	callAi external.CallAI
}

func NewTranslatePage(
	pageRepository *db.PageRepository,
	commentRepository *db.CommentRepository,
	callAi external.CallAI,
) *TranslatePage {
	return &TranslatePage{
		pr:     pageRepository,
		cr:     commentRepository,
		callAi: callAi,
	}
}

type TitleTranslation struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type CommentsForTranslation struct {
	Comments []CommentForTranslation `json:"comments"`
}

type CommentForTranslation struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Score   int    `json:"score"`
}

func NewCommentForTranslation(c db.Comment) CommentForTranslation {
	return CommentForTranslation{
		Id:      c.Id,
		Content: c.Content,
		Score:   -1,
	}
}

func (tp *TranslatePage) run(
	dateString string,
) error {
	pages := tp.pr.FindUntranslatedByDate(dateString)

	for _, page := range pages {
		slog.Info("Translating page", slog.Int("page id", page.Id), slog.String("title", page.Title))

		comments := tp.cr.FindByPageId(page.Id)

		if err := tp.translatePageAndComments(&page, comments); err != nil {
			return err
		}

		// update page, comments
		if err := tp.pr.Update(&page); err != nil {
			slog.Error("failed to update page", "err", err)
			return err
		}

		for _, comment := range comments {
			tp.cr.Update(&comment)
		}
	}

	return nil
}

func (tp *TranslatePage) translatePageAndComments(page *db.Page, comments []db.Comment) error {
	// translate title
	err := tp.retry(func() error {
		return tp.translateTitle(page, comments)
	}, MAX_RETRIES)
	if err != nil {
		return err
	}

	// translate comments
	for i := 0; i < len(comments); i = i + COMMENTS_CHUNK_NUM {
		slog.Info("Translating comment chunk", slog.Int("start index", i))

		commentsChunk := comments[i:min(i+COMMENTS_CHUNK_NUM, len(comments))]

		err = tp.retry(func() error {
			translatedCommentsChunk, err := tp.translateCommentChunk(
				page.TranslatedTitle,
				commentsChunk,
			)
			if err != nil {
				return err
			}

			// update comments
			for j, tc := range translatedCommentsChunk {
				comments[i+j] = tc
			}

			return nil
		}, MAX_RETRIES)

		if err != nil {
			slog.Error("failed to translate comment chunk", "err", err)
			return err
		}

		page.Translated = true
	}

	return nil
}

func (tp *TranslatePage) translateTitle(page *db.Page, comments []db.Comment) error {
	slog.Info("Translating page title", slog.String("title", page.Title))

	commentsForContext := comments[:min(COMMENTS_CONTEXT_NUM, len(comments))]

	// build prompt
	var sb strings.Builder
	for _, comment := range commentsForContext {
		fmt.Fprintln(&sb, comment.Content)
	}

	commentsStr := sb.String()

	prompt := fmt.Sprintf(PROMPT_TITLE, page.Title, commentsStr)

	answer, err := tp.callAi(prompt)
	if err != nil {
		return err
	}

	var titleTranslation TitleTranslation
	if err := json.Unmarshal([]byte(answer), &titleTranslation); err != nil {
		slog.Error("failed to unmarshal", "answer", answer)
		return err
	}

	page.TranslatedTitle = titleTranslation.Title

	page.Tags = strings.Join(titleTranslation.Tags, ",")

	return nil
}

func (tp *TranslatePage) translateCommentChunk(title string, comments []db.Comment) ([]db.Comment, error) {
	commentsForTranslation := make([]CommentForTranslation, 0)
	for _, c := range comments {
		commentsForTranslation = append(
			commentsForTranslation,
			NewCommentForTranslation(c),
		)
	}

	// build a struct for json marshal
	jsonCommentsStruct := CommentsForTranslation{commentsForTranslation}

	jsonComments, err := json.Marshal(jsonCommentsStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %v: %w", jsonCommentsStruct, err)
	}

	prompt := fmt.Sprintf(PROMPT_COMMENT, title, string(jsonComments))

	answer, err := tp.callAi(prompt)
	if err != nil {
		return nil, err
	}

	// For debugging
	// slog.Info(prompt)
	// slog.Info(answer)

	var result CommentsForTranslation

	if err = json.Unmarshal([]byte(answer), &result); err != nil {
		slog.Error("failed to unmarshal", "answer", answer)
		return nil, fmt.Errorf("failed to unmarshal: %s\n %w\n", answer, err)
	}

	if len(result.Comments) != len(comments) {
		slog.Error("Invalid number in answer", "answer", answer)
		return nil, fmt.Errorf("Invalid number of comments: %d != %d, Origina json: %s", len(result.Comments), len(comments), answer)
	}

	for i, translated := range result.Comments {
		comments[i].TranslatedContent = sanitizeTranslatedComment(translated.Content)
		comments[i].Score = translated.Score
		comments[i].Translated = true
	}

	return comments, nil
}

func sanitizeTranslatedComment(c string) string {
	// remove all double quotes
	c = strings.ReplaceAll(c, `"`, `”`)
	// remove all single quotes
	c = strings.ReplaceAll(c, `'`, `’`)

	// replace HTML tags
	// escape <br> first and replace back later
	c = strings.ReplaceAll(c, "\\n", `<br>`)
	c = strings.ReplaceAll(c, "\n", `<br>`)
	c = strings.ReplaceAll(c, `<br>`, `$$br$$`)
	c = strings.ReplaceAll(c, `<`, `＜`)
	c = strings.ReplaceAll(c, `>`, `＞`)
	c = strings.ReplaceAll(c, `$$br$$`, `<br>`)

	return c
}

func (tp *TranslatePage) retry(fn func() error, maxRetries int) error {
	for tried := 0; tried < maxRetries; tried++ {
		err := fn()
		if err != nil {
			slog.Error("running function failed.", "tried", tried+1, "err", err)
			continue
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed to run function after tried %d times", maxRetries)
}
