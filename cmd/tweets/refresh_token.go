package tweets

import (
	"log/slog"

	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/karintomania/kaigai-go-scraper/external"
)

// Renew refresh token and return access token
func getAccessToken(
	kvr *db.KvRepository,
) (string, error) {
	// get current refresh token
	oldRefreshToken := kvr.FindByKey("x_refresh_token")

	// call refresh token endpoint
	res, err := external.RefreshToken(oldRefreshToken)


	if err != nil {
		slog.Error("error on refreshing token", "err", err)
		return "", err
	}

	kvr.Update(&db.Kv{Key: "x_refresh_token", Value: res.RefreshToken})

	//return
	return res.AccessToken, nil
}
