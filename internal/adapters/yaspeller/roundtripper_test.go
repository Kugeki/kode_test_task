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
	rts := make([]*mockRT, 5)
	for i := range rts {
		rts[i] = &mockRT{}
	}

	opts := make([]ClientOpt, 0, len(rts)+1)
	opts = append(opts, WithTimeout(10*time.Second))
	for _, rt := range rts {
		opts = append(opts, WithRoundTripper(rt))
	}

	c, err := NewClient(opts...)
	require.NoError(t, err)

	for _, rt := range rts {
		require.NotNil(t, rt.Next)
	}

	results, err := c.CheckText(context.Background(), "превет")
	assert.NoError(t, err)
	assert.Greater(t, len(results), 0)

	for _, rt := range rts {
		assert.Equal(t, rt.Called, 1)
	}
}
