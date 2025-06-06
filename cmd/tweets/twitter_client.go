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
	accessToken string
}

func NewTwitterClient(kvr *db.KvRepository) *TwitterClient {
	return &TwitterClient{kvr, ""}
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

// TODO: shouldn't be calling this for every tweets
// Renew refresh token and return access token
func (tc *TwitterClient) getAccessToken() (string, error) {
	if tc.accessToken != "" {
		return tc.accessToken, nil
	}

	// get current refresh token
	oldRefreshToken := tc.kvr.FindByKey("x_refresh_token")

	// call refresh token endpoint
	res, err := external.RefreshToken(oldRefreshToken)


	if err != nil {
		slog.Error("error on refreshing token", "err", err)
		return "", err
	}

	tc.kvr.Update(&db.Kv{Key: "x_refresh_token", Value: res.RefreshToken})

	tc.accessToken = res.AccessToken

	//return
	return tc.accessToken, nil
}
