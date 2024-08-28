package yaspeller

import "net/http"

type RoundTripper interface {
	http.RoundTripper

	SetNext(next http.RoundTripper)
}
