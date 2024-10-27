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

type ExportFunc func(w io.Writer, conf *ResumeConfig) error

func exportAndServe(conf *ResumeConfig, f ExportFunc, typ string) http.HandlerFunc {
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

func ExportHTML(w io.Writer, conf *ResumeConfig) error     { return HTMLTemplate.Execute(w, conf) }
func ExportText(w io.Writer, conf *ResumeConfig) error     { return TextTemplate.Execute(w, conf) }
func ExportMarkdown(w io.Writer, conf *ResumeConfig) error { return MarkdownTemplate.Execute(w, conf) }
func ExportJSON(w io.Writer, conf *ResumeConfig) error     { return JSONTemplate.Execute(w, conf) }

func ExportAndServePDF(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportPDF, "application/pdf")
}

func ExportAndServeHTML(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportHTML, "text/html; charset=utf-8")
}

func ExportAndServeJSON(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportJSON, "application/json")
}

func ExportAndServeText(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportText, "text/plain")
}

func ExportAndServeMarkdown(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportMarkdown, "text/markdown")
}
