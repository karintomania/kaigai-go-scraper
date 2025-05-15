package tweets

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestPostScheduledCmd(t *testing.T) {
	common.SetLogger()

	// Mock the database connection and repositories
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	tr := db.NewTweetRepository(dbConn)
	kvr := db.NewKvRepository(dbConn)

	// Insert mock tweets into the database
	mockTweets := []db.Tweet{
		{Id: 1, Content: "Test tweet 1", Published: false, ScheduledAt: "2025-04-27"},
		{Id: 2, Content: "Test tweet 2", Published: false, ScheduledAt: "2025-04-28"},
	}
	for _, tweet := range mockTweets {
		tr.Insert(&tweet)
	}

	// this tweet shouldn't be posted
	futureTweet := db.Tweet{Id: 3, Content: "Test tweet 3", Published: false, ScheduledAt: "2025-04-29"}
	tr.Insert(&futureTweet)

	mockTwitterClient := &MockTwitterClient{
		expectedContents: []string{"Test tweet 1", "Test tweet 2"},
		t:                t,
	}

	// Create command
	cmd := NewTestPostScheduledCmd(tr, kvr, mockTwitterClient)

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

type MockTwitterClient struct {
	expectedContents []string
	t                *testing.T
	counter          int
}

func (tc *MockTwitterClient) Post(content string) error {
	want := tc.expectedContents[tc.counter]

	require.Equal(tc.t, want, content)

	tc.counter++

	return nil
}
