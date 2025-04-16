package external

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"encoding/json"

	"github.com/karintomania/kaigai-go-scraper/common"
)

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	ModelVersion string `json:"modelVersion"`
}

func (gr *GeminiResponse) GetText() string {
	if len(gr.Candidates) == 0 ||
	len(gr.Candidates[0].Content.Parts) == 0 {
		log.Panicf("Invalid gemini response: %v", gr)
	}

	return gr.Candidates[0].Content.Parts[0].Text
}

func CallGemini(prompt string) string {

	// data, _ := os.ReadFile("./external/sample_responses/gemini.json")
	data := geminiHttpCall(prompt)

	var gr GeminiResponse

	if err := json.Unmarshal(data, &gr); err != nil {
		log.Fatalln(err)
	}

	answer := gr.GetText()

	return answer
}

func geminiHttpCall(prompt string) []byte {
	url := fmt.Sprintf(
		common.GetEnv("gemini_url"),
		common.GetEnv("gemini_model"),
		common.GetEnv("gemini_api_key"),
		)

	body := []byte(fmt.Sprintf(`{
    "contents": [
      {
        "parts": [
          {
            "text": "%s"
          }
        ]
      }
    ]
}`, prompt))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		slog.Error("failed to create request to gemini", slog.Any("err", err))
		panic(err)
	}

	// TODO: reuse http client
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http call to gemini failed", slog.Any("err", err))
		panic(err)
	}

	b := resp.Body
	defer b.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("gemini returned error status", slog.String("status", resp.Status))
		panic(resp.Status)
	}

	responseBytes, err := io.ReadAll(b)
	if err != nil {
		slog.Error("failed to read gemini response body", slog.Any("err", err))
		panic(err)
	}

	return responseBytes
}
