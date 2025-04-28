package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTweetRepository(t *testing.T) {
	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	tr := NewTweetRepository(db)

	tr.CreateTable()

	t.Run("Insert", func(t *testing.T) {
		tr.Truncate()

		tweet := Tweet{
			PageId:      1,
			Date:        "2025-01-01",
			Content:     "Test Tweet",
			ScheduledAt: "2025-01-02",
			Published:   false,
		}

		tr.Insert(&tweet)

		tweets := tr.FindUnpublishedByScheduledDate("2025-01-02")

		require.Equal(t, 1, len(tweets), "Expected %d tweet, got %d", 1, len(tweets))

		created := tweets[0]

		require.Equal(t, tweet.PageId, created.PageId)
		require.Equal(t, tweet.Date, created.Date)
		require.Equal(t, tweet.Content, created.Content)
		require.Equal(t, tweet.ScheduledAt, created.ScheduledAt)
		require.Equal(t, tweet.Published, created.Published)
	})

	t.Run("Update", func(t *testing.T) {
		tr.Truncate()

		// Insert initial tweet
		tweet := Tweet{
			PageId:      1,
			Date:        "2025-01-01",
			Content:     "Test Tweet",
			ScheduledAt: "2025-01-02",
			Published:   false,
		}
		tr.Insert(&tweet)

		// Update the tweet
		updatedTweet := &Tweet{
			Id:          tweet.Id,
			PageId:      2,
			Date:        "2025-01-03",
			Content:     "Updated Tweet",
			ScheduledAt: "2025-01-04",
			Published:   true,
		}
		tr.Update(updatedTweet)

		// Verify the update
		tweets := tr.FindUnpublishedByScheduledDate("2025-01-04")
		require.Equal(t, 0, len(tweets), "Expected %d tweet, got %d", 0, len(tweets))

		tweets = tr.FindUnpublishedByScheduledDate("2025-01-05")
		require.Equal(t, 0, len(tweets), "Expected %d tweet, got %d", 0, len(tweets))
	})

	t.Run("FindUnpublishedByScheduledDate", func(t *testing.T) {
		tr.Truncate()

		tweets := []Tweet{
			// Scheduled on the day
			{
				PageId:      1,
				Content:     "Test Tweet 1",
				ScheduledAt: "2025-01-05",
				Published:   false,
			},
			// Scheduled on the day before
			{
				PageId:      2,
				Content:     "Test Tweet 2",
				ScheduledAt: "2025-01-01",
				Published:   false,
			},
			// Future scheduled tweet
			{
				PageId:      3,
				Content:     "Test Tweet 3",
				ScheduledAt: "2025-02-31",
				Published:   false,
			},
			// Already published
			{
				PageId:      4,
				Content:     "Test Tweet 4",
				ScheduledAt: "2025-01-02",
				Published:   true,
			},
		}

		for _, tweet := range tweets {
			tr.Insert(&tweet)
		}

		result := tr.FindUnpublishedByScheduledDate("2025-01-05")

		require.Equal(t, 2, len(result))

		require.Equal(t, 1, result[0].PageId)
		require.Equal(t, "Test Tweet 1", result[0].Content)
		require.Equal(t, 2, result[1].PageId)
		require.Equal(t, "Test Tweet 2", result[1].Content)
	})

	t.Run("FindUnpublishedByScheduledDate returns nothing for no record", func(t *testing.T) {
		tr.Truncate()

		tweets := tr.FindUnpublishedByScheduledDate("2025-01-01")

		require.Equal(t, 0, len(tweets), "Expected %d tweet, got %d", 0, len(tweets))
	})
}
