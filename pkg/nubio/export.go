package nubio

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

type ExportType string

const (
	ExportTypeHTML ExportType = "html"
	ExportTypePDF  ExportType = "pdf"
	ExportTypeJSON ExportType = "json"
)

var tmplFuncs = template.FuncMap{
	"subtract": func(a, b int) int { return a - b },
}

func mustParseHTMLTmpl(name, raw string) *template.Template {
	return template.Must(template.New(name).Funcs(tmplFuncs).Parse(raw))
}

var (
	//go:embed resume.html.gotmpl
	HTMLRawTemplate string
	HTMLTemplate    = mustParseHTMLTmpl("html", HTMLRawTemplate)
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

func ExportHTML(w io.Writer, conf *ResumeConfig) error { return HTMLTemplate.Execute(w, conf) }

func ExportJSON(w io.Writer, conf *ResumeConfig) error {
	return json.NewEncoder(w).Encode(conf.ToResumeExport())
}

func ExportAndServePDF(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportPDF, "application/pdf")
}

func ExportAndServeHTML(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportHTML, "text/html; charset=utf-8")
}

func ExportAndServeJSON(conf *ResumeConfig) http.HandlerFunc {
	return exportAndServe(conf, ExportJSON, "application/json")
}
