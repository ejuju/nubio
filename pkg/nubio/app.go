package nubio

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ejuju/nubio/pkg/httpmux"
)

const (
	PathHome        = "/"
	PathFaviconSVG  = "/favicon.svg"
	PathSitemapXML  = "/sitemap.xml"
	PathRobotsTXT   = "/robots.txt"
	PathProfileJSON = "/profile.json"
	PathProfilePDF  = "/profile.pdf"
	PathProfileTXT  = "/profile.txt"
	PathProfileMD   = "/profile.md"
	PathPGPKey      = "/pgp.asc"
)

func RunApp(args ...string) (exitcode int) {
	// Init logger.
	slogh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogh)
	logger.Debug("logger ready")

	// Load config.
	configPath := "local.config.json"
	if len(args) > 0 {
		configPath = args[0]
	}
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("read config", "error", err)
		return 1
	}
	config := &Config{}
	err = json.Unmarshal(rawConfig, config)
	if err != nil {
		logger.Error("parse config", "error", err)
		return 1
	}

	// Load user profile.
	rawProfile, err := os.ReadFile(config.Profile)
	if err != nil {
		logger.Error("read profile", "error", err)
		return 1
	}
	profile := &Profile{}
	err = json.Unmarshal(rawProfile, profile)
	if err != nil {
		logger.Error("parse profile", "error", err)
		return 1
	}

	// Init and register HTTP endpoints.
	endpoints := httpmux.Map{
		PathHome:        {"GET": ExportAndServeHTML(profile)},
		PathFaviconSVG:  {"GET": serveFaviconSVG()},
		PathSitemapXML:  {"GET": serveSitemapXML(profile.Domain)},
		PathRobotsTXT:   {"GET": serveRobotsTXT()},
		PathProfileJSON: {"GET": ExportAndServeJSON(profile)},
		PathProfilePDF:  {"GET": ExportAndServePDF(profile)},
		PathProfileTXT:  {"GET": ExportAndServeText(profile)},
		PathProfileMD:   {"GET": ExportAndServeMarkdown(profile)},
	}

	// If provided, load PGP key and register endpoint to serve it.
	if profile.Contact.PGP != "" {
		pgpKey, err := os.ReadFile(profile.Contact.PGP)
		if err != nil {
			logger.Error("read PGP public key file", "error", err)
			return 1
		}
		endpoints[PathPGPKey] = map[string]http.Handler{"GET": servePGPKey(pgpKey)}
	}

	router := endpoints.Handler(http.NotFoundHandler())

	// Wrap global middleware.
	router = httpmux.Wrap(router,
		httpmux.NewPanicRecoveryMiddleware(handlePanic(logger)),
	)

	// Run HTTP server.
	// TODO: Setup TLS and use HTTPS in prod.
	s := &http.Server{
		Addr:              config.Address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    50_000,
	}
	errc := make(chan error, 1)
	go func() { errc <- s.ListenAndServe() }()

	// Wait for interrupt or server error.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	select {
	case err := <-errc:
		logger.Error("critical failure", "error", err)
		return
	case sig := <-interrupt:
		logger.Debug("shutting down", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		logger.Error("shutdown HTTP server", "error", err)
	}

	logger.Debug("shutdown successful")
	return 0
}

func handlePanic(logger *slog.Logger) httpmux.PanicRecoveryHandler {
	return func(w http.ResponseWriter, r *http.Request, err any) {
		logger.Error("handler panicked", "error", err, "path", r.URL.Path, "address", r.RemoteAddr)
	}
}

func serveFaviconSVG() http.HandlerFunc {
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

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, content)
	}
}

func serveRobotsTXT() http.HandlerFunc {
	const content = `User-Agent: *
Allow: /
`

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, content)
	}
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
