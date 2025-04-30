package db

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/stretchr/testify/require"
)

func TestCommentRepository(t *testing.T) {
	common.SetLogger()

	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	cr := NewCommentRepository(db)
	cr.CreateTable()

	t.Run("Insert", func(t *testing.T) {
		cr.Truncate()
		comment := Comment{
			ExtCommentId:      "cmt123",
			PageId:            1,
			UserName:          "test_user",
			Content:           "This is a test comment.",
			TranslatedContent: "これはテストコメントです。",
			Indent:            0,
			Reply:             0,
			Colour:            "blue",
			Score:             10,
			Translated:        true,
		}

		cr.Insert(&comment)

		comments := cr.FindByPageId(comment.Id)
		require.Equal(t, 1, len(comments), "Expected %d comment, got %d", 1, len(comments))

		created := comments[0]
		require.Equal(t, comment.ExtCommentId, created.ExtCommentId)
		require.Equal(t, comment.PageId, created.PageId)
		require.Equal(t, comment.UserName, created.UserName)
		require.Equal(t, comment.Content, created.Content)
		require.Equal(t, comment.TranslatedContent, created.TranslatedContent)
		require.Equal(t, comment.Indent, created.Indent)
		require.Equal(t, comment.Reply, created.Reply)
		require.Equal(t, comment.Colour, created.Colour)
		require.Equal(t, comment.Score, created.Score)
		require.Equal(t, comment.Translated, created.Translated)
	})

	t.Run("FindByPageIdReturnNothing", func(t *testing.T) {
		cr.Truncate()
		comments := cr.FindByPageId(1)
		require.Equal(t, 0, len(comments), "Expected %d comment, got %d", 0, len(comments))
	})

	t.Run("Update", func(t *testing.T) {
		cr.Truncate()
		pageId := 1

		// Insert initial comment
		comment := Comment{
			ExtCommentId:      "cmt123",
			PageId:            pageId,
			UserName:          "test_user",
			Content:           "This is a test comment.",
			TranslatedContent: "これはテストコメントです。",
			Indent:            0,
			Reply:             0,
			Colour:            "blue",
			Score:             10,
		}
		cr.Insert(&comment)

		// Update the comment
		updatedComment := Comment{
			Id:                comment.Id,
			ExtCommentId:      "cmt456",
			PageId:            pageId,
			UserName:          "updated_user",
			Content:           "This is an updated comment.",
			TranslatedContent: "これは更新されたコメントです。",
			Indent:            1,
			Reply:             0,
			Colour:            "red",
			Score:             20,
		}
		cr.Update(&updatedComment)

		// Verify the update
		comments := cr.FindByPageId(pageId)
		require.Equal(t, 1, len(comments), "Expected %d comment, got %d", 1, len(comments))

		updated := comments[0]
		require.Equal(t, updatedComment.ExtCommentId, updated.ExtCommentId)
		require.Equal(t, updatedComment.PageId, updated.PageId)
		require.Equal(t, updatedComment.UserName, updated.UserName)
		require.Equal(t, updatedComment.Content, updated.Content)
		require.Equal(t, updatedComment.TranslatedContent, updated.TranslatedContent)
		require.Equal(t, updatedComment.Indent, updated.Indent)
		require.Equal(t, updatedComment.Reply, updated.Reply)
		require.Equal(t, updatedComment.Colour, updated.Colour)
		require.Equal(t, updatedComment.Score, updated.Score)
	})

	t.Run("DoesExtIdExist", func(t *testing.T) {
		cr.Truncate()

		// Insert initial comment
		comment := Comment{
			ExtCommentId: "cmt_1",
			PageId:       1,
		}
		cr.Insert(&comment)

		require.True(t, cr.DoesExtIdExist(1, "cmt_1"))
		require.False(t, cr.DoesExtIdExist(1, "cmt_2"))
	})
}
