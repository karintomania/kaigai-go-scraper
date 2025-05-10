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

func Post() {
	url := "https://api.x.com/2/tweets"

	title := "会社への忠誠心ってマジ意味ある？ 企業に尽くしても報われない現実（2018年）"
	link := "https://www.kaigai-tech-matome.com/posts/2025_04/2025_04_24_on_loyalty_to_your_employer_2018/"
	content := fmt.Sprintf("「%s」に対する海外の反応をまとめました。 #テックニュース #海外の反応\n\n%s", title, link)

	payload := strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, content))

	req, _ := http.NewRequest("POST", url, payload)

	bearer := ""
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

func RefreshToken(refreshToken string) (RefreshTokenResponse, error) {
	response := RefreshTokenResponse{}
	baseUrl := common.GetEnv("x_base_url")
	tokenUrl := fmt.Sprintf("%s/2/oauth2/token", baseUrl)

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", common.GetEnv("x_client_id"))
	formData.Set("refresh_token", refreshToken)

	slog.Info("debug form data", "formdata", formData.Encode())
	payload := strings.NewReader(formData.Encode())

	req, _ := http.NewRequest("POST", tokenUrl, payload)

	bearer := common.GetEnv("x_bearer_token")
	req.Method = "POST"
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	cli := getHttpClient()

	res, _ := cli.Do(req)

	defer res.Body.Close()

	var body []byte
	body, err := io.ReadAll(res.Body)

	if err != nil {
		slog.Error("failed to read response body", "err", err)
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
