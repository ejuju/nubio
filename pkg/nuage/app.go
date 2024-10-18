package nuage

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ejuju/nubio/pkg/httpmux"
	"github.com/ejuju/nubio/pkg/nubio"
)

type Config struct {
	Address      string `json:"address"`        // Local HTTP server address.
	Profile      string `json:"profile"`        // Path to JSON file where profile data is stored.
	TrueIPHeader string `json:"true_ip_header"` // Ex: "X-Forwarded-For", useful when reverse proxying.
	PGPKey       string `json:"pgp_key"`        // Path to PGP public key file.
}

func Run(args ...string) (exitcode int) {
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
	profile := &nubio.Profile{}
	err = json.Unmarshal(rawProfile, profile)
	if err != nil {
		logger.Error("parse profile", "error", err)
		return 1
	}

	// Load PGP key if provided.
	var pgpKey []byte
	if config.PGPKey != "" {
		pgpKey, err = os.ReadFile(config.PGPKey)
		if err != nil {
			logger.Error("read PGP public key file", "error", err)
			return 1
		}
		profile.Contact.PGP = profile.Domain + PathPGPKey
	}

	// Init and register HTTP endpoints.
	endpoints := httpmux.Map{
		PathPing:        {"GET": http.HandlerFunc(servePing)},
		PathFaviconSVG:  {"GET": http.HandlerFunc(serveFaviconSVG)},
		PathRobotsTXT:   {"GET": http.HandlerFunc(serveRobotsTXT)},
		PathSitemapXML:  {"GET": serveSitemapXML(profile.Domain)},
		PathHome:        {"GET": nubio.ExportAndServeHTML(profile)},
		PathProfileJSON: {"GET": nubio.ExportAndServeJSON(profile)},
		PathProfilePDF:  {"GET": nubio.ExportAndServePDF(profile)},
		PathProfileTXT:  {"GET": nubio.ExportAndServeText(profile)},
		PathProfileMD:   {"GET": nubio.ExportAndServeMarkdown(profile)},
		PathPGPKey:      {"GET": servePGPKey(pgpKey)},
	}

	router := endpoints.Handler(http.NotFoundHandler())

	// Wrap global middleware.
	//
	// Note: the panic recovery middleware relies on the:
	//	- True IP middleware
	//	- Request ID middleware
	//	- Logging middleware (to know if a response has been sent).
	//
	// This also means that any panic occuring in one of the above mentioned
	// middlewares propagates up and will cause the program to exit.
	router = httpmux.Wrap(router,
		httpmux.NewTrueIPMiddleware(config.TrueIPHeader),
		httpmux.NewRequestIDMiddleware(),
		httpmux.NewLoggingMiddleware(handleAccessLog(logger)),
		httpmux.NewPanicRecoveryMiddleware(handlePanic(logger)),
	)

	// Run HTTP server in separate Goroutine.
	// TODO: Support HTTPS.
	httpServerErrLogger := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	httpServerErrLogger.SetPrefix("http server: ")
	s := &http.Server{
		Addr:              config.Address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    50_000,
		ErrorLog:          httpServerErrLogger,
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

	// Shutdown HTTP server gracefully.
	err = s.Shutdown(ctx)
	if err != nil {
		logger.Error("shutdown HTTP server", "error", err)
	}

	// Done.
	logger.Debug("shutdown successful")
	return 0
}
