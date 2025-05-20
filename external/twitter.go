package external

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/karintomania/kaigai-go-scraper/common"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Post(content string, bearer string) error {
	baseUrl := common.GetEnv("x_base_url")
	url := fmt.Sprintf("%s/2/tweets", baseUrl)

	payload := strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, content))

	slog.Info("pyaload", "payload", payload)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))
	req.Header.Add("Content-Type", "application/json")

	cli := getHttpClient()

	res, _ := cli.Do(req)
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post tweet: %s", res.Status)
	}
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}

func RefreshToken(refreshToken string) (RefreshTokenResponse, error) {
	response := RefreshTokenResponse{}
	baseUrl := common.GetEnv("x_base_url")
	tokenUrl := fmt.Sprintf("%s/2/oauth2/token", baseUrl)

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", common.GetEnv("x_client_id"))
	formData.Set("refresh_token", refreshToken)

	payload := strings.NewReader(formData.Encode())

	req, _ := http.NewRequest("POST", tokenUrl, payload)

	bearer := common.GetEnv("x_bearer_token")
	req.Method = "POST"
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", bearer))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	cli := getHttpClient()

	res, _ := cli.Do(req)

	defer res.Body.Close()

	var body []byte
	body, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		errMsg := "error on refresh token request"
		slog.Error(errMsg, "status", res.StatusCode, "body", string(body))
		return response, fmt.Errorf(errMsg)
	}

	if err != nil {
		slog.Error("failed to read response body", "err", err, "body", string(body))
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("failed to unmarshal response body", "body", body, "err", err)
		return response, err
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return response, nil
}
