package nubio

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

var (
	//go:embed resume.html.gotmpl
	HTMLRawTemplate string
	HTMLTemplate    = template.Must(template.New("html").Parse(HTMLRawTemplate))

	//go:embed resume.txt.gotmpl
	TextRawTemplate string
	TextTemplate    = template.Must(template.New("txt").Parse(TextRawTemplate))

	//go:embed resume.md.gotmpl
	MarkdownRawTemplate string
	MarkdownTemplate    = template.Must(template.New("md").Parse(MarkdownRawTemplate))
)

type ExportFunc func(w io.Writer, p *Config) error

func ExportAndServe(conf *Config, f ExportFunc, typ string) http.HandlerFunc {
	buf := &bytes.Buffer{}
	err := f(buf, conf)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", typ)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func ExportHTML(w io.Writer, conf *Config) error     { return HTMLTemplate.Execute(w, conf) }
func ExportText(w io.Writer, conf *Config) error     { return TextTemplate.Execute(w, conf) }
func ExportMarkdown(w io.Writer, conf *Config) error { return MarkdownTemplate.Execute(w, conf) }
func ExportJSON(w io.Writer, conf *Config) error     { return json.NewEncoder(w).Encode(conf.Resume) }

func ExportAndServePDF(conf *Config) http.HandlerFunc {
	return ExportAndServe(conf, ExportPDF, "application/pdf")
}

func ExportAndServeHTML(conf *Config) http.HandlerFunc {
	return ExportAndServe(conf, ExportHTML, "text/html; charset=utf-8")
}

func ExportAndServeJSON(conf *Config) http.HandlerFunc {
	return ExportAndServe(conf, ExportJSON, "application/json")
}

func ExportAndServeText(conf *Config) http.HandlerFunc {
	return ExportAndServe(conf, ExportText, "text/plain")
}

func ExportAndServeMarkdown(conf *Config) http.HandlerFunc {
	return ExportAndServe(conf, ExportMarkdown, "text/markdown")
}
