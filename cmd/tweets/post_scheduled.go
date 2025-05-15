package tweets

import (
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

type PostScheduledCmd struct {
	tr  *db.TweetRepository
	kvr *db.KvRepository
	tc  Poster
}

func NewPostScheduledCmd() *PostScheduledCmd {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))
	defer dbConn.Close()

	tr := db.NewTweetRepository(dbConn)
	kvr := db.NewKvRepository(dbConn)

	return &PostScheduledCmd{
		tr:  tr,
		kvr: kvr,
		tc:  NewTwitterClient(kvr),
	}
}

func NewTestPostScheduledCmd(
	tr *db.TweetRepository,
	kvr *db.KvRepository,
	tc Poster,
) *PostScheduledCmd {
	return &PostScheduledCmd{
		tr:  tr,
		kvr: kvr,
		tc:  tc,
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

		err := cmd.tc.Post(tweet.Content)
		if err != nil {
			slog.Error("error on posting tweet", "tweet_id", tweet.Id, "error", err)
			return err
		}

		tweet.Published = true
		cmd.tr.Update(&tweet)
	}

	return nil
}
