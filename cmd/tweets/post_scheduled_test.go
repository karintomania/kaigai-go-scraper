package tweets

import (
	"fmt"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestPostScheduledCmd(t *testing.T) {
	// Mock the database connection and repositories
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	tr := db.NewTweetRepository(dbConn)

	// Insert mock tweets into the database
	mockTweets := []db.Tweet{
		{Id: 1, Content: "Test tweet 1", Published: false, Date: "2025-04-28"},
		{Id: 2, Content: "Test tweet 2", Published: false, Date: "2025-04-28"},
	}
	for _, tweet := range mockTweets {
		tr.Insert(&tweet)
	}

	counter := 1
	mockPostTweet := func(tweetContent string) error {
		require.Equal(t, fmt.Sprintf("Test tweet %d", counter), tweetContent)
		counter++
		return nil
	}

	// Create command
	cmd := NewTestPostScheduledCmd(tr, mockPostTweet)

	// Execute command
	err := cmd.Run("2025-04-28")
	require.NoError(t, err)

	// Verify that the tweets were marked as published
	for _, tweet := range mockTweets {
		tweetFromDb := tr.FindById(tweet.Id)
		require.NoError(t, err)
		require.NotNil(t, tweetFromDb)

		require.True(t, tweetFromDb.Published)
	}
}
