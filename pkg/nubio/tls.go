package nubio

import (
	"log/slog"
	"net/http"

	"github.com/ejuju/nubio/pkg/httpmux"
	"golang.org/x/crypto/acme/autocert"
)

func listenAndServeProd(h http.Handler, config *Config, logger *slog.Logger) {
	// Configure autocert.
	tlsCertManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domain, "www."+config.Domain),
		Cache:      autocert.DirCache(config.TLSDirpath),
		Email:      config.TLSEmailAddress,
	}

	// Serve autocert on port :80.
	autocertServer := httpmux.NewDefaultHTTPServer(":80", tlsCertManager.HTTPHandler(nil), logger)
	autocertServerErrChan := make(chan error, 1)
	go func() {
		err := autocertServer.ListenAndServe()
		if err != nil {
			autocertServerErrChan <- err
		}
	}()

	// Serve app on port :443.
	httpsServer := httpmux.NewDefaultHTTPServer(":443", h, logger)
	httpsServerErrChan := make(chan error, 1)
	go func() {
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			httpsServerErrChan <- err
		}
	}()
}
