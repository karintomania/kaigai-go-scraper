package external

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActuallyCallGemini(t *testing.T) {
	// This test is only for local testing
	return

	prompt := "just reply 'test'"

	result, err := CallGemini(prompt)
	require.NoError(t, err)

	require.Equal(t, result, "Test.\n")
}

func TestEscapeStringForJson(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"n\n", "n\\n"},
		{"t\t", "t\\t"},
		{"slash\\slash", "slash\\\\slash"},
		{"\"quoted\"", "quoted"},
	}
	for _, tc := range testCases {
		result := escapeStringForJSON(tc.input)
		require.Equal(t, tc.expected, result)
	}
}
