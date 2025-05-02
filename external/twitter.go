package external

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
