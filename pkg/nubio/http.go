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
	PathPing        = "/ping"
	PathFaviconSVG  = "/favicon.svg"
	PathSitemapXML  = "/sitemap.xml"
	PathRobotsTXT   = "/robots.txt"
	PathProfileHTML = "/"
	PathProfileJSON = "/profile.json"
	PathProfilePDF  = "/profile.pdf"
	PathProfileTXT  = "/profile.txt"
	PathProfileMD   = "/profile.md"
	PathPGPKey      = "/pgp.asc"
)

func NewHTTPHandler(fallback http.Handler, profile *Profile, pgpKey []byte) http.Handler {
	return httpmux.Map{
		PathPing:        {"GET": http.HandlerFunc(servePing)},
		PathFaviconSVG:  {"GET": http.HandlerFunc(serveFaviconSVG)},
		PathRobotsTXT:   {"GET": http.HandlerFunc(serveRobotsTXT)},
		PathSitemapXML:  {"GET": serveSitemapXML(profile.Domain)},
		PathProfileHTML: {"GET": ExportAndServeHTML(profile)},
		PathProfilePDF:  {"GET": ExportAndServePDF(profile)},
		PathProfileJSON: {"GET": ExportAndServeJSON(profile)},
		PathProfileTXT:  {"GET": ExportAndServeText(profile)},
		PathProfileMD:   {"GET": ExportAndServeMarkdown(profile)},
		PathPGPKey:      {"GET": servePGPKey(pgpKey)},
	}.Handler(fallback)
}

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

		// Log error.
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
	io.WriteString(w, "pong\n")
}

//go:embed favicon.svg
var faviconSVG []byte

func serveFaviconSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(faviconSVG)
}

const robotsTXT = `User-Agent: *
Disallow:
`

func serveRobotsTXT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, robotsTXT)
}

func generateSitemapXML(domain string) []byte {
	paths := []string{
		PathProfileHTML,
		PathProfileJSON,
		PathProfilePDF,
		PathProfileTXT,
		PathProfileMD,
	}

	b := &bytes.Buffer{}
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, path := range paths {
		b.WriteString("<url><loc>https://" + domain + path + "/</loc></url>\n")
	}
	b.WriteString("</urlset>\n")

	return b.Bytes()
}

func serveSitemapXML(domain string) http.HandlerFunc {
	content := generateSitemapXML(domain)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}

func servePGPKey(key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(key)
	}
}
