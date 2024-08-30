package yaspeller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/Kugeki/kode_test_task/pkg/reqlog"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

type LogRoundTripper struct {
	log    *slog.Logger
	reqLog *reqlog.RequestLogger

	next http.RoundTripper

	requestIDPrefix string
	requestCounter  atomic.Uint32
}

var _ RoundTripper = (*LogRoundTripper)(nil)

var RequestIDHeader = "X-Request-Id"

func NewLogRoundTripper(log *slog.Logger, defaultLevel slog.Level) *LogRoundTripper {
	var buf [12]byte

	rand.Read(buf[:])
	prefix := base64.URLEncoding.EncodeToString(buf[:])

	return &LogRoundTripper{
		log:             log,
		reqLog:          reqlog.New(log, defaultLevel),
		requestIDPrefix: prefix,
	}
}

func (l *LogRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestID := r.Header.Get(RequestIDHeader)
	if requestID == "" {
		requestID = fmt.Sprintf("%s-%06d", l.requestIDPrefix, l.requestCounter.Add(1))
	}

	l.reqLog.LogBegin(r.Context(), r, requestID)

	t1 := time.Now()
	resp, err := l.next.RoundTrip(r)

	l.reqLog.LogEnd(r.Context(), requestID, resp.StatusCode, time.Since(t1))
	return resp, err
}

func (l *LogRoundTripper) SetNext(next http.RoundTripper) {
	l.next = next
}
