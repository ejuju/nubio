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
		logger.Error("missing arguments", "args", []string{"path to profile.json", "path to output directory"})
		return 1
	}
	profilePath := args[0]
	outputDirpath := args[1]

	// Load profile.
	profile, err := LoadProfileFile(profilePath)
	if err != nil {
		logger.Error("load profile config", "error", err)
		return 1
	}

	// List export paths and corresponding function.
	exports := map[string]ExportFunc{
		"index.html":                             ExportHTML,
		strings.TrimPrefix(PathProfilePDF, "/"):  ExportPDF,
		strings.TrimPrefix(PathProfileJSON, "/"): ExportJSON,
		strings.TrimPrefix(PathProfileTXT, "/"):  ExportText,
		strings.TrimPrefix(PathProfileMD, "/"):   ExportMarkdown,
	}

	// Generate static files.
	//
	// Note that PGP key is not included here
	// since we only rely on the profile.json,
	// and the local PGP key is specified in the server.json.
	files := map[string][]byte{
		strings.TrimPrefix(PathFaviconSVG, "/"): faviconSVG,
		strings.TrimPrefix(PathSitemapXML, "/"): generateSitemapXML(profile.Domain),
		strings.TrimPrefix(PathRobotsTXT, "/"):  []byte(robotsTXT),
		strings.TrimPrefix(PathPing, "/"):       []byte("ok\n"),
		strings.TrimPrefix(PathVersion, "/"):    []byte(version + "\n"),
	}
	for path, export := range exports {
		b := &bytes.Buffer{}
		err = export(b, profile)
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
