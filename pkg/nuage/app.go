package nuage

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ejuju/nuage/pkg/httpmux"
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

	// Generate exports.

	// Init and register HTTP endpoints.
	endpoints := httpmux.Map{
		"/":             {"GET": ExportAndServeHTML(profile)},
		"/favicon.ico":  {"GET": http.NotFoundHandler()},
		"/sitemap.xml":  {"GET": http.NotFoundHandler()},
		"/robots.txt":   {"GET": http.NotFoundHandler()},
		"/profile.md":   {"GET": ExportAndServeMarkdown(profile)},
		"/profile.json": {"GET": ExportAndServeJSON(profile)},
		"/profile.txt":  {"GET": ExportAndServeText(profile)},
		"/profile.pdf":  {"GET": ExportAndServePDF(profile)},
	}
	router := endpoints.Handler(http.NotFoundHandler())

	// TODO: Wrap global middleware.
	// - Panic recovery
	// - Logging
	// - IP ban (rate limiting + blocklist)
	// - Auth

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
