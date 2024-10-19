package nubio

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ejuju/nubio/pkg/httpmux"
	"golang.org/x/crypto/acme/autocert"
)

type Config struct {
	Address         string `json:"address"`        // Local HTTP server address.
	Domain          string `json:"domain"`         // Public domain name used to host the site.
	TrueIPHeader    string `json:"true_ip_header"` // Ex: "X-Forwarded-For", useful when reverse proxying.
	TLSDirpath      string `json:"tls_dirpath"`    // Path to TLS certificate directory.
	TLSEmailAddress string `json:"tls_email_addr"` // Email address in TLS certificate.
	Profile         string `json:"profile"`        // Path to JSON file where profile data is stored.
	PGPKey          string `json:"pgp_key"`        // Path to PGP public key file.
}

func Run(args ...string) (exitcode int) {
	// Init logger.
	slogh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogh)
	logger.Debug("logger ready")

	// Load config.
	configPath := "config.json"
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
	profile.Domain = config.Domain

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
		PathProfileHTML: {"GET": ExportAndServeHTML(profile)},
		PathProfilePDF:  {"GET": ExportAndServePDF(profile)},
		PathProfileJSON: {"GET": ExportAndServeJSON(profile)},
		PathProfileTXT:  {"GET": ExportAndServeText(profile)},
		PathProfileMD:   {"GET": ExportAndServeMarkdown(profile)},
		PathPGPKey:      {"GET": servePGPKey(pgpKey)},
	}

	router := endpoints.Handler(http.NotFoundHandler())

	// Wrap global middleware.
	//
	// Note: the panic recovery middleware relies on the:
	//	- True IP middleware (for debugging purposes).
	//	- Request ID middleware (for debugging purposes).
	//	- Logging middleware (to respond to the client if the panic occured before write).
	//
	// This also means that any panic occuring in one of the above mentioned
	// middlewares propagates up and will cause the program to exit.
	//
	// Other middlewares should be put below the panic recovery middleware.
	router = httpmux.Wrap(router,
		httpmux.NewTrueIPMiddleware(config.TrueIPHeader),
		httpmux.NewRequestIDMiddleware(),
		httpmux.NewLoggingMiddleware(handleAccessLog(logger)),
		httpmux.NewPanicRecoveryMiddleware(handlePanic(logger)),
		httpmux.RedirectToNonWWW,
	)

	// Run HTTP(S) server(s).
	if config.TLSDirpath != "" {
		exitcode = runHTTPS(router, config, logger)
	} else {
		exitcode = runHTTP(router, config, logger)
	}

	// Done.
	logger.Info("exiting", "code", exitcode)
	return exitcode
}

func runHTTP(h http.Handler, config *Config, logger *slog.Logger) (exitcode int) {
	errc := make(chan error, 1)
	s := httpmux.NewDefaultHTTPServer(config.Address, h, logger)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			errc <- err
		}
	}()

	// Wait for interrupt or server error.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	select {
	case err := <-errc:
		logger.Error("critical failure", "error", err)
		return 1
	case sig := <-interrupt:
		logger.Debug("shutting down", "signal", sig.String())
	}

	// Shutdown HTTP server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		logger.Error("shutdown HTTP server", "error", err)
		exitcode = 1
	}

	return exitcode
}

func runHTTPS(h http.Handler, config *Config, logger *slog.Logger) (exitcode int) {
	// Ensure certs directory exists.
	fstat, err := os.Stat(config.TLSDirpath)
	if err != nil {
		logger.Error("check certs directory", "error", err)
		return 1
	} else if !fstat.IsDir() {
		logger.Error("check certs directory", "error", "not a directory")
		return 1
	}

	// Configure autocert.
	tlsCertManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domain, "www."+config.Domain),
		Cache:      autocert.DirCache(config.TLSDirpath),
		Email:      config.TLSEmailAddress,
	}
	autocertServer := httpmux.NewDefaultHTTPServer(":80", tlsCertManager.HTTPHandler(nil), logger)
	httpsServer := httpmux.NewDefaultHTTPServer(":443", h, logger)
	httpsServer.TLSConfig = tlsCertManager.TLSConfig()
	errc := make(chan error, 1)

	// Serve autocert on port :80.
	go func() {
		err := autocertServer.ListenAndServe()
		if err != nil {
			errc <- fmt.Errorf("run HTTP server: %w", err)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			httpsServer.Shutdown(ctx)
		}
	}()

	// Serve app on port :443.
	go func() {
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			errc <- fmt.Errorf("run HTTPS server: %w", err)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			autocertServer.Shutdown(ctx)
		}
	}()

	// Wait for interrupt or server error.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	select {
	case err := <-errc:
		logger.Error("critical failure", "error", err)
		return 1
	case sig := <-interrupt:
		logger.Debug("shutting down", "signal", sig.String())
	}

	// Shutdown HTTP server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = autocertServer.Shutdown(ctx)
	if err != nil {
		logger.Error("shutdown HTTP server", "error", err)
		exitcode = 1
	}

	// Shutdown HTTPS server gracefully.
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = httpsServer.Shutdown(ctx)
	if err != nil {
		logger.Error("shutdown HTTP server", "error", err)
		exitcode = 1
	}

	return exitcode
}
