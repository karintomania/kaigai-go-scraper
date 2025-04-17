package db

import (
	"testing"
)

func TestCommentRepositoryInsert(t *testing.T) {
	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	cr := NewCommentRepository(db)

	cr.CreateTable()

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

	if want, got := 1, len(comments); want != got {
		t.Fatalf("Expected %d comment, got %d", want, got)
	}

	created := comments[0]

	if created.ExtCommentId != comment.ExtCommentId ||
		created.PageId != comment.PageId ||
		created.UserName != comment.UserName ||
		created.Content != comment.Content ||
		created.TranslatedContent != comment.TranslatedContent ||
		created.Indent != comment.Indent ||
		created.Reply != comment.Reply ||
		created.Colour != comment.Colour ||
		created.Score != comment.Score ||
		created.Translated != comment.Translated {
		t.Fatalf("Expected %v, got %v", comment, created)
	}
}

func TestCommentRepositoryFindByPageIdReturnNothing(t *testing.T) {
	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	cr := NewCommentRepository(db)

	cr.CreateTable()

	comments := cr.FindByPageId(1)

	if want, got := 0, len(comments); want != got {
		t.Fatalf("Expected %d comment, got %d", want, got)
	}
}

func TestCommentRepositoryUpdate(t *testing.T) {
	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	cr := NewCommentRepository(db)
	cr.CreateTable()

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
	if want, got := 1, len(comments); want != got {
		t.Fatalf("Expected %d comment, got %d", want, got)
	}

	updated := comments[0]
	if updated.ExtCommentId != updatedComment.ExtCommentId ||
		updated.PageId != updatedComment.PageId ||
		updated.UserName != updatedComment.UserName ||
		updated.Content != updatedComment.Content ||
		updated.TranslatedContent != updatedComment.TranslatedContent ||
		updated.Indent != updatedComment.Indent ||
		updated.Reply != updatedComment.Reply ||
		updated.Colour != updatedComment.Colour ||
		updated.Score != updatedComment.Score {
		t.Fatalf("Expected %v, got %v", updatedComment, updated)
	}
}
