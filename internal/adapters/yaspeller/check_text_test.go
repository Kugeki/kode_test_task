package yaspeller

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_CheckText(t *testing.T) {
	c, err := NewClient(WithTimeout(10 * time.Second))
	require.NoError(t, err)

	results, err := c.CheckText(context.Background(), "теставое зодание")
	assert.NoError(t, err)
	assert.Greater(t, len(results), 0)
}
