package nubio

import (
	"bytes"
	_ "embed"
	"io"
	"log/slog"
	"net/http"

	"github.com/ejuju/nubio/pkg/httpmux"
)

const (
	PathHome        = "/"
	PathPing        = "/ping"
	PathFaviconSVG  = "/favicon.svg"
	PathSitemapXML  = "/sitemap.xml"
	PathRobotsTXT   = "/robots.txt"
	PathProfileJSON = "/profile.json"
	PathProfilePDF  = "/profile.pdf"
	PathProfileTXT  = "/profile.txt"
	PathProfileMD   = "/profile.md"
	PathPGPKey      = "/pgp.asc"
)

func handleAccessLog(logger *slog.Logger) httpmux.LoggingHandlerFunc {
	return func(w *httpmux.ResponseRecorderWriter, r *http.Request) {
		logger.Info("handled HTTP",
			"status", w.StatusCode,
			"path", r.URL.Path,
			"written", w.Written,
			"request_id", httpmux.GetRequestID(r.Context()),
			"ip_address", httpmux.GetTrueIP(r.Context()),
		)
	}
}

func handlePanic(logger *slog.Logger) httpmux.PanicRecoveryHandler {
	return func(w http.ResponseWriter, r *http.Request, err any) {
		// Write response to client if none has been written yet.
		rrw, ok := w.(*httpmux.ResponseRecorderWriter)
		panicBeforeResponse := ok && rrw.StatusCode == -1
		if panicBeforeResponse {
			http.Error(w, "Critical failure", http.StatusInternalServerError)
		}

		// Logger error.
		logger.Error("handler panicked",
			"error", err,
			"path", r.URL.Path,
			"request_id", httpmux.GetRequestID(r.Context()),
		)
	}
}

func servePing(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "pong")
}

//go:embed favicon.svg
var faviconSVG []byte

func serveFaviconSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(faviconSVG)
}

func serveRobotsTXT(w http.ResponseWriter, r *http.Request) {
	const content = `User-Agent: *
Disallow:
`
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, content)
}

func serveSitemapXML(domain string) http.HandlerFunc {
	paths := []string{
		PathHome,
		PathProfileJSON,
		PathProfilePDF,
		PathProfileTXT,
		PathProfileMD,
	}

	b := &bytes.Buffer{}
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")
	for _, path := range paths {
		b.WriteString("<url><loc>https://" + domain + path + "/</loc></url>")
	}
	b.WriteString("</urlset>")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(b.Bytes())
	}
}

func servePGPKey(key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(key)
	}
}
