package nubio

import (
	"bytes"
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

func serveFaviconSVG(w http.ResponseWriter, r *http.Request) {
	const content = `<svg xmlns="http://www.w3.org/2000/svg" version="1.1" viewBox="0 0 13 11" shape-rendering="crispEdges">
	<rect width="1" height="11" x="0" y="0" fill="#cacaca" />
	<rect width="1" height="11" x="12" y="0" fill="#cacaca" />
	<rect width="3" height="1" x="1" y="0" fill="#cacaca" />
	<rect width="3" height="1" x="9" y="0" fill="#cacaca" />
	<rect width="3" height="1" x="9" y="10" fill="#cacaca" />
	<rect width="3" height="1" x="1" y="10" fill="#cacaca" />
	<rect width="2" height="1" x="3" y="2" fill="#cacaca" />
	<rect width="2" height="1" x="8" y="2" fill="#cacaca" />
	<rect width="4" height="1" x="2" y="3" fill="#cacaca" />
	<rect width="4" height="1" x="7" y="3" fill="#cacaca" />
	<rect width="9" height="1" x="2" y="4" fill="#cacaca" />
	<rect width="7" height="1" x="3" y="5" fill="#cacaca" />
	<rect width="5" height="1" x="4" y="6" fill="#cacaca" />
	<rect width="3" height="1" x="5" y="7" fill="#cacaca" />
	<rect width="1" height="1" x="6" y="8" fill="#cacaca" />
</svg>
`
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, content)
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
