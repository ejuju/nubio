package nubio

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func RunSSG(args ...string) (exitcode int) {
	slogh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogh)

	// Check arguments.
	if len(args) < 2 {
		logger.Error("missing arguments", "args", []string{"path to config.json", "path to output directory"})
		return 1
	}
	configPath := args[0]
	outputDirpath := args[1]

	// Load conf.
	conf, err := LoadResumeConfig(configPath)
	if err != nil {
		logger.Error("load resume config", "error", err)
		return 1
	}
	errs := conf.Check()
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error("bad resume config", "error", err)
		}
		return 1
	}

	// List export paths and corresponding function.
	exports := map[string]ExportFunc{
		"index.html":                            ExportHTML,
		strings.TrimPrefix(PathResumePDF, "/"):  ExportPDF,
		strings.TrimPrefix(PathResumeJSON, "/"): ExportJSON,
		strings.TrimPrefix(PathResumeTXT, "/"):  ExportText,
		strings.TrimPrefix(PathResumeMD, "/"):   ExportMarkdown,
	}

	// Generate static files.
	files := map[string][]byte{
		strings.TrimPrefix(PathFaviconSVG, "/"): faviconSVG,
		strings.TrimPrefix(PathSitemapXML, "/"): generateSitemapXML(conf.Domain),
		strings.TrimPrefix(PathRobotsTXT, "/"):  []byte(robotsTXT),
		strings.TrimPrefix(PathPing, "/"):       []byte("ok\n"),
		strings.TrimPrefix(PathVersion, "/"):    []byte(version + "\n"),
	}
	if len(conf.PGPKey) > 0 {
		files[strings.TrimPrefix(PathPGPKey, "/")] = []byte(conf.PGPKey)
	}
	if len(conf.CustomCSS) > 0 {
		files[strings.TrimPrefix(PathCustomCSS, "/")] = []byte(conf.CustomCSS)
	}
	for path, export := range exports {
		b := &bytes.Buffer{}
		err = export(b, conf)
		if err != nil {
			logger.Error("export", "path", path, "error", err)
			return 1
		}
		files[path] = b.Bytes()
	}

	// Write files.
	for path, f := range files {
		path = filepath.Join(outputDirpath, path)
		err := os.WriteFile(path, f, 0666)
		if err != nil {
			logger.Error("write file", "path", path, "error", err)
			return 1
		}
		logger.Info("wrote file", "path", path)
	}

	logger.Info("all files written")
	return 0
}
