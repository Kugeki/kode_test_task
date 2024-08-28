package yaspeller

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client

	url string
}

var DefaultURL = "https://speller.yandex.net/services/spellservice.json"

func NewClient(opts ...ClientOpt) (*Client, error) {
	c := &Client{}

	c.client = &http.Client{}
	c.url = DefaultURL

	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.MaxIdleConnsPerHost = 100

	c.client.Transport = defaultTransport

	for _, op := range opts {
		err := op(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) post(ctx context.Context, path string, contentType string, body io.Reader) (*http.Response, error) {
	u, err := url.JoinPath(c.url, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
