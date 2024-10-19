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
	//go:embed profile.html.gotmpl
	HTMLRawTemplate string
	HTMLTemplate    = template.Must(template.New("html").Parse(HTMLRawTemplate))

	//go:embed profile.txt.gotmpl
	TextRawTemplate string
	TextTemplate    = template.Must(template.New("txt").Parse(TextRawTemplate))

	//go:embed profile.md.gotml
	MarkdownRawTemplate string
	MarkdownTemplate    = template.Must(template.New("md").Parse(MarkdownRawTemplate))
)

type ExportFunc func(w io.Writer, p *Profile) error

func ExportAndServe(p *Profile, f ExportFunc, typ string) http.HandlerFunc {
	buf := &bytes.Buffer{}
	err := f(buf, p)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", typ)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func ExportHTML(w io.Writer, p *Profile) error     { return HTMLTemplate.Execute(w, p) }
func ExportJSON(w io.Writer, p *Profile) error     { return json.NewEncoder(w).Encode(p) }
func ExportText(w io.Writer, p *Profile) error     { return TextTemplate.Execute(w, p) }
func ExportMarkdown(w io.Writer, p *Profile) error { return MarkdownTemplate.Execute(w, p) }
