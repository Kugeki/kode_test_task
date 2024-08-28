package yaspeller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type mockRT struct {
	Called int
	Next   http.RoundTripper
}

func (rt *mockRT) SetNext(next http.RoundTripper) {
	rt.Next = next
}

func (rt *mockRT) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.Called += 1

	return rt.Next.RoundTrip(r)
}

func TestWithRoundTripper(t *testing.T) {
	rt := &mockRT{}

	c, err := NewClient(WithTimeout(10*time.Second), WithRoundTripper(rt))
	require.NoError(t, err)

	require.NotNil(t, rt.Next)

	results, err := c.CheckText(context.Background(), "превет", "", "")
	assert.NoError(t, err)
	assert.Greater(t, len(results), 0)

	assert.Equal(t, rt.Called, 1)

}
