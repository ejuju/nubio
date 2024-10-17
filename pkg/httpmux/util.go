package httpmux

import "net/http"

type ResponseRecorderWriter struct {
	http.ResponseWriter
	StatusCode int
	Written    int
}

var _ http.ResponseWriter = (*ResponseRecorderWriter)(nil)

func (w *ResponseRecorderWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *ResponseRecorderWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.Written += n
	return n, err
}
