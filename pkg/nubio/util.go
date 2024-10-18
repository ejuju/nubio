package nubio

import "net/http"

func ExportAndServeMarkdown(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportMarkdown, "text/markdown")
}

func ExportAndServeText(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportText, "text/plain")
}

func ExportAndServeJSON(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportJSON, "application/json")
}

func ExportAndServeHTML(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportHTML, "text/html; charset=utf-8")
}

func ExportAndServePDF(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportPDF, "application/pdf")
}
