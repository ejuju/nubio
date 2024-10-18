package httpmux

import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseRecorderWriter struct {
	http.ResponseWriter
	StatusCode int
	Written    int
}

var _ http.ResponseWriter = (*ResponseRecorderWriter)(nil)

func (w *ResponseRecorderWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *ResponseRecorderWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.Written += n
	return n, err
}

// Returns a HTTP server configured with sensible timeouts and limits.
func NewDefaultHTTPServer(addr string, h http.Handler, logger *slog.Logger) *http.Server {
	errLogger := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	errLogger.SetPrefix(addr + ": ")
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    100_000,
		ErrorLog:          errLogger,
	}
}
