package external

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
	t.Skip("This test is only for E2E test")
	Post()
	require.True(t, true)
}
