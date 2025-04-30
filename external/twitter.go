package external

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Post() {

	url := "https://api.twitter.com/2/tweets"

	payload := strings.NewReader(`{
	"text": "Learn how to use the user Tweet timeline and user mention timeline endpoints in the X API v2 to explore Tweetu"
}`)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Authorization", "Bearer <token>")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
