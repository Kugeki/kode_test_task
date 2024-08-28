package yaspeller

import (
	"net/http"
	"time"
)

type ClientOpt func(c *Client) error

func WithRoundTripper(rt RoundTripper) ClientOpt {
	return func(c *Client) error {
		transport := c.client.Transport
		rt.SetNext(transport)
		c.client.Transport = rt
		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) error {
		c.client.Timeout = timeout
		return nil
	}
}

func WithHost(host string) ClientOpt {
	return func(c *Client) error {
		c.url = host
		return nil
	}
}

func WithTransport(transport http.RoundTripper) ClientOpt {
	return func(c *Client) error {
		c.client.Transport = transport
		return nil
	}
}

const (
	IgnoreDigits         = 2
	IgnoreURLs           = 4
	FindRepeatWords      = 8
	IgnoreCapitalization = 512
)

type CheckOpt func(old int) int

func WithCheckIgnoreDigits() CheckOpt {
	return func(old int) int {
		return old | IgnoreDigits
	}
}

func WithCheckIgnoreURLs() CheckOpt {
	return func(old int) int {
		return old | IgnoreURLs
	}
}

func WithCheckFindRepeatWords() CheckOpt {
	return func(old int) int {
		return old | FindRepeatWords
	}
}

func WithCheckIgnoreCapitalization() CheckOpt {
	return func(old int) int {
		return old | IgnoreCapitalization
	}
}
