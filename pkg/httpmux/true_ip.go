package httpmux

import (
	"context"
	"net"
	"net/http"
)

type trueIPCtxType string

const trueIPCtxKey trueIPCtxType = "true-IP"

func WithTrueIP(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, trueIPCtxKey, v)
}

func GetTrueIP(ctx context.Context) string {
	v, ok := ctx.Value(trueIPCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func NewTrueIPMiddleware(xForwardedFor string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get remote address from http.Server,
			// or from header if provided (used for reverse proxying).
			v := r.RemoteAddr
			if xForwardedFor != "" {
				v = r.Header.Get(xForwardedFor)
			}

			// Try to get the IP address,
			// otherwise, just leave original remote adddress (or header value).
			ipAddr, _, err := net.SplitHostPort(v)
			if err == nil {
				v = ipAddr
			}

			h.ServeHTTP(w, r.WithContext(WithTrueIP(r.Context(), v)))
		})
	}
}
