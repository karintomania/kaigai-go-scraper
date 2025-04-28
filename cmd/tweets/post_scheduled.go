package tweets

import (
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/db"
)

type PostScheduledCmd struct {
	tr        *db.TweetRepository
	postTweet func(tweetContent string) error
}

// func NewPostScheduledCmd( tr *db.TweetRepository) *PostScheduledCmd {
// 	return &PostScheduledCmd{
// 		tr:         tr,
// 		postTweet:  postTweet,
// 	}
// }

func NewTestPostScheduledCmd(
	tr *db.TweetRepository,
	postTweet func(tweetContent string) error,
) *PostScheduledCmd {
	return &PostScheduledCmd{
		tr:        tr,
		postTweet: postTweet,
	}
}

func (cmd *PostScheduledCmd) Run(dateStr string) error {
	tweets := cmd.tr.FindUnpublishedByScheduledDate(dateStr)

	if len(tweets) == 0 {
		slog.Info("No tweets scheduled for today", "date", dateStr)
		return nil
	}

	for _, tweet := range tweets {
		slog.Info("Posting tweet", "tweet_id", tweet.Id)

		err := cmd.postTweet(tweet.Content)
		if err != nil {
			slog.Error("error on posting tweet", "tweet_id", tweet.Id, "error", err)
			return err
		}

		tweet.Published = true
		cmd.tr.Update(&tweet)
	}

	return nil
}
