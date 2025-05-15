package tweets

import (
	"strconv"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestScheduleTweets(t *testing.T) {
	// Mock the database connection and repositories
	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	pr := db.NewPageRepository(dbConn)
	tr := db.NewTweetRepository(dbConn)

	// Create a ScheduleTweetsCmd instance
	cmd := ScheduleTweetsCmd{
		pr: pr,
		tr: tr,
	}

	// Define the test date
	testDate := "2025-10-01"
	scheduledDateMax := "2025-10-08"

	page := db.Page{
		Slug:  "test_slug",
		Title: "Test Title",
		Date:  testDate,
	}

	pr.Insert(&page)

	// Run the command
	err := cmd.Run(testDate, []int{page.Id})
	require.NoError(t, err)

	// Verify that the expected pages were processed
	tweets := tr.FindUnpublishedByScheduledDate(scheduledDateMax)

	require.Equal(t, page.Id, tweets[0].PageId)
	require.Equal(t, testDate, tweets[0].Date)

	expectedContent := `「Test Title」に対する海外の反応をまとめました。
#海外の反応 #テックニュース

https://www.kaigai-tech-matome.com/posts/2025_10/test_slug/`
	require.Equal(t, expectedContent, tweets[0].Content)
	require.Equal(t, false, tweets[0].Published)

	scheduleDay, err := strconv.Atoi(tweets[0].ScheduledAt[9:])
	require.NoError(t, err)

	require.GreaterOrEqual(t, scheduleDay, 2)
	require.LessOrEqual(t, scheduleDay, 8)
}
