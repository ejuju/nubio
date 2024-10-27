package nubio

import (
	"bytes"
	_ "embed"
	templatehtml "html/template"
	"io"
	"net/http"
	templatetxt "text/template"
)

type ExportFormat string

const (
	ExportTypeHTML ExportFormat = "html"
	ExportTypePDF  ExportFormat = "pdf"
	ExportTypeJSON ExportFormat = "json"
	ExportTypeTXT  ExportFormat = "txt"
	ExportTypeMD   ExportFormat = "md"
)

var tmplFuncs = templatehtml.FuncMap{
	"subtract": func(a, b int) int { return a - b },
}

func mustParseHTMLTmpl(name, raw string) *templatehtml.Template {
	return templatehtml.Must(templatehtml.New(name).Funcs(tmplFuncs).Parse(raw))
}

func mustParseTextTmpl(name, raw string) *templatetxt.Template {
	return templatetxt.Must(templatetxt.New(name).Funcs(tmplFuncs).Parse(raw))
}

var (
	//go:embed resume.html.gotmpl
	HTMLRawTemplate string
	HTMLTemplate    = mustParseHTMLTmpl("html", HTMLRawTemplate)

	//go:embed resume.txt.gotmpl
	TextRawTemplate string
	TextTemplate    = mustParseTextTmpl("txt", TextRawTemplate)

	//go:embed resume.json.gotmpl
	JSONRawTemplate string
	JSONTemplate    = mustParseTextTmpl("json", JSONRawTemplate)

	//go:embed resume.md.gotmpl
	MarkdownRawTemplate string
	MarkdownTemplate    = mustParseTextTmpl("md", MarkdownRawTemplate)
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
func ExportJSON(w io.Writer, conf *Config) error     { return JSONTemplate.Execute(w, conf) }

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
