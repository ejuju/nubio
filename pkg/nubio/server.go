package nubio

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ejuju/nubio/pkg/httpmux"
	"golang.org/x/crypto/acme/autocert"
)

type ServerConfig struct {
	Address         string `json:"address"`        // Local HTTP server address.
	TrueIPHeader    string `json:"true_ip_header"` // Optional: (ex: "X-Forwarded-For", use when reverse proxying).
	TLSDirpath      string `json:"tls_dirpath"`    // Path to TLS certificate directory.
	TLSEmailAddress string `json:"tls_email_addr"` // Email address in TLS certificate.
	ResumePath      string `json:"resume_path"`    // Path to resume config file.
}

// Read and decode server and resume config files.
func LoadServerAndResumeConfig(path string) (*ServerConfig, *ResumeConfig, error) {
	serverConf, err := LoadServerConfig(path)
	if err != nil {
		return nil, nil, fmt.Errorf("load server config: %w", err)
	}

	resumeConf, err := LoadResumeConfig(serverConf.ResumePath)
	if err != nil {
		return nil, nil, fmt.Errorf("load resume config: %w", err)
	}

	return serverConf, resumeConf, nil
}

func LoadServerConfig(path string) (conf *ServerConfig, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	conf = &ServerConfig{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	return conf, nil
}

func RunServer(args ...string) (exitcode int) {
	// Init logger.
	slogh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogh)
	logger.Debug("logger ready")

	// Load config.
	defaultServerConfigPath := "server.json"
	if len(args) > 0 {
		defaultServerConfigPath = args[0]
	}
	serverConf, resumeConf, err := LoadServerAndResumeConfig(defaultServerConfigPath)
	if err != nil {
		logger.Error("load config", "error", err)
		return 1
	}
	errs := serverConf.Check()
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error("bad server config", "error", err)
		}
		return 1
	}
	errs = resumeConf.Check()
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error("bad resume config", "error", err)
		}
		return 1
	}

	// Init and register HTTP endpoints, wrap global middleware.
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
	h := httpmux.Wrap(NewHTTPHandler(nil, resumeConf),
		httpmux.NewTrueIPMiddleware(serverConf.TrueIPHeader),
		httpmux.NewRequestIDMiddleware(),
		httpmux.NewLoggingMiddleware(handleAccessLog(logger)),
		httpmux.NewPanicRecoveryMiddleware(handlePanic(logger)),
		httpmux.RedirectToNonWWW,
	)

	// Run HTTP(S) server(s).
	if serverConf.TLSDirpath != "" {
		exitcode = runHTTPS(h, serverConf, resumeConf, logger)
	} else {
		exitcode = runHTTP(h, serverConf, logger)
	}

	// Done.
	logger.Info("exiting", "code", exitcode)
	return exitcode
}

func runHTTP(h http.Handler, config *ServerConfig, logger *slog.Logger) (exitcode int) {
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
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
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

func runHTTPS(h http.Handler, serverConf *ServerConfig, resumeConf *ResumeConfig, logger *slog.Logger) (exitcode int) {
	// Ensure certs directory exists.
	fstat, err := os.Stat(serverConf.TLSDirpath)
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
		HostPolicy: autocert.HostWhitelist(resumeConf.Domain, "www."+resumeConf.Domain),
		Cache:      autocert.DirCache(serverConf.TLSDirpath),
		Email:      serverConf.TLSEmailAddress,
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
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
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
