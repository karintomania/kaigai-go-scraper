package external

import (
	"testing"
)

func TestCallGemini(t *testing.T) {
	// This test is only for local testing
	return

	prompt := "just reply 'test'"

	result := CallGemini(prompt)
	if result != "Test.\n" {
		t.Errorf("Expected answer, got '%s'", result)
	}
}
