package httpmux

import (
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

// Wraps middlewares around given handler (ordered from outermost to innermost).
func Wrap(h http.Handler, v ...Middleware) http.Handler {
	for i := len(v) - 1; i >= 0; i-- {
		h = v[i](h)
	}
	return h
}

type PanicRecoveryHandler func(w http.ResponseWriter, r *http.Request, err any)

func NewPanicRecoveryMiddleware(onPanic PanicRecoveryHandler) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					onPanic(w, r, err)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
