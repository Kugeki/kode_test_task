package yaspeller

import (
	"log/slog"
	"net/http"
	"time"
)

type ClientOpt func(c *Client) error

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

func WithLogger(log *slog.Logger, defaultLevel slog.Level) ClientOpt {
	return WithRoundTripper(
		NewLogRoundTripper(log.With(slog.String("context", "yaspeller client")), defaultLevel),
	)
}

func WithRoundTripper(rt RoundTripper) ClientOpt {
	return func(c *Client) error {
		transport := c.client.Transport
		rt.SetNext(transport)
		c.client.Transport = rt
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

type CheckConfig struct {
	Lang    string
	Format  string
	Options int
}

type CheckOpt func(cfg *CheckConfig)

func WithCheckLang(lang string) CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Lang = lang
	}
}

func WithCheckFormat(format string) CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Format = format
	}
}

func WithCheckIgnoreDigits() CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Options |= IgnoreDigits
	}
}

func WithCheckIgnoreURLs() CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Options |= IgnoreURLs
	}
}

func WithCheckFindRepeatWords() CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Options |= FindRepeatWords
	}
}

func WithCheckIgnoreCapitalization() CheckOpt {
	return func(cfg *CheckConfig) {
		cfg.Options |= IgnoreCapitalization
	}
}
