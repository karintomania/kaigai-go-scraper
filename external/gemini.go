package external

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/karintomania/kaigai-go-scraper/common"
)

type CallAI func(string) (string, error)

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

func (gr *geminiResponse) getText() (string, error) {
	if len(gr.Candidates) == 0 ||
		len(gr.Candidates[0].Content.Parts) == 0 {
		slog.Error("Invalid gemini response", "gr", gr)
		return "", fmt.Errorf("Invalid gemini response: %v", gr)
	}

	return gr.Candidates[0].Content.Parts[0].Text, nil
}

func CallGemini(prompt string) (string, error) {
	sleepSecond := time.Duration(common.GetEnvInt("gemini_sleep_second"))
	time.Sleep(sleepSecond * time.Second)

	var data []byte
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		data, err = geminiHttpCall(prompt)
		wg.Done()
	}()

	// This basically wait at least the specifed second, and then as long as gemini takes
	time.Sleep(sleepSecond * time.Second)
	wg.Wait()

	if err != nil {
		return "", err
	}

	var gr geminiResponse

	if err := json.Unmarshal(data, &gr); err != nil {
		return "", err
	}

	answer, err := gr.getText()
	if err != nil {
		return "", err
	}

	answer = sanitizeResponse(answer)

	return answer, nil
}

func geminiHttpCall(prompt string) ([]byte, error) {
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
		slog.Error("failed to create request to gemini", "err", err)
		return nil, err
	}

	client := getHttpClient()

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("http call to gemini failed", "err", err)
		return nil, err
	}

	b := resp.Body
	defer b.Close()

	responseBytes, err := io.ReadAll(b)
	if err != nil {
		slog.Error("failed to read gemini response body", "err", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("gemini returned error status", "status", resp.Status, "body", string(responseBytes))
		return nil, fmt.Errorf("gemini returned error status: %s", resp.Status)
	}

	return responseBytes, nil
}

// Use json.Marshal to make the string json compatible
func escapeStringForJSON(s string) string {
	escapedBytes, err := json.Marshal(s)
	if err != nil {
		slog.Error("failed to escape the string", "s", s)
		panic(err)
	}

	escapedString := string(escapedBytes)
	// Remove surrounding quotes
	escapedString = strings.Trim(escapedString, "\\\"")
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
