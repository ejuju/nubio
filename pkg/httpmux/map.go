package httpmux

import "net/http"

type Map map[string]map[string]http.Handler

func (m Map) Handler(fallback http.Handler) http.Handler {
	if fallback == nil {
		fallback = http.NotFoundHandler()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers, ok := m[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		handler, ok := handlers[r.Method]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
