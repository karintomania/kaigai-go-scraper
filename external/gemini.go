package external

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/karintomania/kaigai-go-scraper/common"
)

type CallAI func(string) string

// dogde the rate limit
const GEMINI_SLEEP_SECONDS = 5

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	ModelVersion string `json:"modelVersion"`
}

func (gr *geminiResponse) getText() string {
	if len(gr.Candidates) == 0 ||
		len(gr.Candidates[0].Content.Parts) == 0 {
		log.Panicf("Invalid gemini response: %v", gr)
	}

	return gr.Candidates[0].Content.Parts[0].Text
}

func CallGemini(prompt string) string {
	time.Sleep(GEMINI_SLEEP_SECONDS * time.Second)

	data := geminiHttpCall(prompt)

	var gr geminiResponse

	if err := json.Unmarshal(data, &gr); err != nil {
		log.Fatalln(err)
	}

	answer := gr.getText()

	answer = sanitizeResponse(answer)

	slog.Info("gemini answer", slog.String("answer", answer))

	return answer
}

func geminiHttpCall(prompt string) []byte {
	url := fmt.Sprintf(
		common.GetEnv("gemini_url"),
		common.GetEnv("gemini_model"),
		common.GetEnv("gemini_api_key"),
	)

	body := []byte(fmt.Sprintf(`{"contents": [
{"parts": [{"text": "%s"}]}
]}`, escapeStringForJSON(prompt)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		slog.Error("failed to create request to gemini", slog.Any("err", err))
		panic(err)
	}

	client := getHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http call to gemini failed", slog.Any("err", err))
		panic(err)
	}

	b := resp.Body
	defer b.Close()

	responseBytes, err := io.ReadAll(b)
	if err != nil {
		slog.Error("failed to read gemini response body", slog.Any("err", err))
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("gemini returned error status", slog.String("status", resp.Status), slog.String("body", string(responseBytes)))
		panic(resp.Status)
	}

	return responseBytes
}

func escapeStringForJSON(s string) string {
	escapedBytes, err := json.Marshal([]string{s})
	if err != nil {
		log.Panic(err)
	}
	// Remove [" and "]
	escapedString := string(escapedBytes[2 : len(escapedBytes)-2])

	return escapedString
}

// this removes code block markers (```json or ```) from the beginning and end of a string.
// Gemini often add these quotes.
func sanitizeResponse(answer string) string {
	answer = strings.TrimSpace(answer)

	answer = strings.TrimPrefix(answer, "```json")
	answer = strings.TrimPrefix(answer, "```")

	answer = strings.TrimSuffix(answer, "```")

	return strings.TrimSpace(answer)
}
