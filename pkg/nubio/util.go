package nubio

import "net/http"

func ExportAndServePDF(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportPDF, "application/pdf")
}

func ExportAndServeHTML(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportHTML, "text/html; charset=utf-8")
}

func ExportAndServeJSON(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportJSON, "application/json")
}

func ExportAndServeText(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportText, "text/plain")
}

func ExportAndServeMarkdown(p *Profile) http.HandlerFunc {
	return ExportAndServe(p, ExportMarkdown, "text/markdown")
}
