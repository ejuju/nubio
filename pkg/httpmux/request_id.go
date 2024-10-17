package httpmux

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
)

type requestIDCtxType string

const requestIDCtxKey requestIDCtxType = "true-IP"

func WithRequestID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, v)
}

func GetRequestID(ctx context.Context) string {
	v, ok := ctx.Value(requestIDCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func NewRequestIDMiddleware() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b := [8]byte{}
			_, err := io.ReadFull(rand.Reader, b[:])
			v := ""
			if err == nil {
				v = hex.EncodeToString(b[:])
			}
			h.ServeHTTP(w, r.WithContext(WithRequestID(r.Context(), v)))
		})
	}
}
