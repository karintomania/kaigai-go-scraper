package external

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
)

func TestCallGemini(t *testing.T) {
	// This test is only for local testing
	return

	prompt := "just reply 'test'"
	key := "test"

	cfg := common.Config{
		"gemini_api_key": key,
		"gemini_url": "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		"gemini_model": "gemini-2.0-flash",
	}

	result := CallGemini(prompt, cfg)
	if result != "Test.\n" {
		t.Errorf("Expected answer, got '%s'", result)
	}
}
