package httpmux

import (
	"log/slog"
	"net/http"
	"strings"
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

// Intercepts requests on "www." subdomain and redirects them to the non-www equivalent.
// Otherwise, just forwards request to the next handler.
func RedirectToNonWWW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, "www.") {
			scheme := "https://"
			if r.TLS == nil {
				scheme = "http://"
			}
			url := scheme + strings.TrimPrefix(r.Host, "www.") + r.URL.Path
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		}

		h.ServeHTTP(w, r)
	})
}
