package tweets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	kvr := db.NewKvRepository(dbConn)
	kvr.Insert(&db.Kv{Key: "x_refresh_token", Value: "refresh_token"})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())

		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/2/oauth2/token", r.URL.String())
		require.Equal(t, "Bearer bearer_token", r.Header.Get("Authorization"))
		require.Equal(t, "client_id", r.PostFormValue("client_id"))
		require.Equal(t, "refresh_token", r.FormValue("refresh_token"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(
			`{"access_token": "access_token", "refresh_token": "new_refresh_token"}`,
		))
	}))
	defer s.Close()

	common.MockEnv("x_client_id", "client_id")
	common.MockEnv("x_bearer_token", "bearer_token")
	common.MockEnv("x_base_url", s.URL)

	res, err := getAccessToken(kvr)
	require.NoError(t, err)

	require.Equal(t, "access_token", res)
}
