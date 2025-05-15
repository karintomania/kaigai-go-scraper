package tweets

import (
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

type Poster interface {
	Post(content string) error
}

type TwitterClient struct {
	kvr *db.KvRepository
}

func NewTwitterClient(kvr *db.KvRepository) *TwitterClient {
	return &TwitterClient{kvr}
}

func (tc *TwitterClient) Post(content string) error {
	// get current access token
	accessToken, err := tc.getAccessToken()
	if err != nil {
		return err
	}

	// call post tweet endpoint
	err = external.Post(content, accessToken)
	if err != nil {
		slog.Error("error on posting tweet", "err", err)
		return err
	}

	return nil
}

// Renew refresh token and return access token
func (tc *TwitterClient) getAccessToken() (string, error) {
	// get current refresh token
	oldRefreshToken := tc.kvr.FindByKey("x_refresh_token")

	// call refresh token endpoint
	res, err := external.RefreshToken(oldRefreshToken)

	if err != nil {
		slog.Error("error on refreshing token", "err", err)
		return "", err
	}

	tc.kvr.Update(&db.Kv{Key: "x_refresh_token", Value: res.RefreshToken})

	//return
	return res.AccessToken, nil
}
