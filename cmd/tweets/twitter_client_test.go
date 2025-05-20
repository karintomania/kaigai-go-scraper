package tweets

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
	"github.com/stretchr/testify/require"
)

func TestTwitterClient(t *testing.T) {
	common.SetLogger()

	dbConn, cleanup := db.GetTestDbConnection()
	defer cleanup()

	kvr := db.NewKvRepository(dbConn)
	kvr.Insert(&db.Kv{Key: "x_refresh_token", Value: "refresh_token"})

	requestCount := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount == 0 {
			assertRefreshTokenRequest(t, w, r)
			requestCount++
		} else {
			assertPostRequest(t, w, r)
		}
	}))
	defer s.Close()

	common.MockEnv("x_client_id", "client_id")
	common.MockEnv("x_bearer_token", "bearer_token")
	common.MockEnv("x_base_url", s.URL)

	tc := NewTwitterClient(kvr)
	err := tc.Post("test")
	require.NoError(t, err)

	// assert updated token
	updatedToken := kvr.FindByKey("x_refresh_token")
	require.Equal(t, "new_refresh_token", updatedToken)
}

func assertRefreshTokenRequest(t *testing.T, w http.ResponseWriter, r *http.Request) {
	require.NoError(t, r.ParseForm())
	require.Equal(t, "POST", r.Method)
	require.Equal(t, "/2/oauth2/token", r.URL.String())
	require.Equal(t, "Basic bearer_token", r.Header.Get("Authorization"))
	require.Equal(t, "client_id", r.PostFormValue("client_id"))
	require.Equal(t, "refresh_token", r.FormValue("refresh_token"))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		`{"access_token": "access_token", "refresh_token": "new_refresh_token"}`,
	))
}

func assertPostRequest(t *testing.T, w http.ResponseWriter, r *http.Request) {
	require.Equal(t, "POST", r.Method)
	require.Equal(t, "/2/tweets", r.URL.String())
	require.Equal(t, "Bearer access_token", r.Header.Get("Authorization"))

	body, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	require.Equal(t, `{"text":"test"}`, string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(
		`Success`,
	))
}
