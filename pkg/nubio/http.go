package nubio

import (
	"bytes"
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/ejuju/nubio/pkg/httpmux"
)

const (
	PathPing       = "/ping"
	PathVersion    = "/version"
	PathFaviconSVG = "/favicon.svg"
	PathSitemapXML = "/sitemap.xml"
	PathRobotsTXT  = "/robots.txt"
	PathResumeHTML = "/"
	PathResumeJSON = "/resume.json"
	PathResumePDF  = "/resume.pdf"
	PathResumeTXT  = "/resume.txt"
	PathResumeMD   = "/resume.md"
	PathPGPKey     = "/pgp.asc"
	PathCustomCSS  = "/custom.css"
)

func NewHTTPHandler(fallback http.Handler, conf *ResumeConfig) http.Handler {
	m := httpmux.Map{
		PathPing:       {"GET": httpmux.TextHandler("ok\n")},
		PathVersion:    {"GET": httpmux.TextHandler(version + "\n")},
		PathFaviconSVG: {"GET": httpmux.SVGHandler(faviconSVG)},
		PathRobotsTXT:  {"GET": httpmux.TextHandler(robotsTXT)},
		PathSitemapXML: {"GET": httpmux.XMLHandler(generateSitemapXML(conf.Domain))},
		PathResumeHTML: {"GET": ExportAndServeHTML(conf)},
		PathResumePDF:  {"GET": ExportAndServePDF(conf)},
		PathResumeJSON: {"GET": ExportAndServeJSON(conf)},
		PathResumeTXT:  {"GET": ExportAndServeText(conf)},
		PathResumeMD:   {"GET": ExportAndServeMarkdown(conf)},
	}
	if len(conf.PGPKey) > 0 {
		m[PathPGPKey] = map[string]http.Handler{"GET": httpmux.TextHandler(string(conf.PGPKey))}
	}
	if len(conf.CustomCSS) > 0 {
		m[PathCustomCSS] = map[string]http.Handler{"GET": httpmux.CSSHandler([]byte(conf.CustomCSS))}
	}

	return m.Handler(fallback)
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

//go:embed favicon.svg
var faviconSVG []byte

const robotsTXT = `User-Agent: *
Disallow:
`

func generateSitemapXML(domain string) []byte {
	paths := []string{
		PathResumeHTML,
		PathResumeJSON,
		PathResumePDF,
		PathResumeTXT,
		PathResumeMD,
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
