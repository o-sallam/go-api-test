package utils

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipHandler wraps an http.HandlerFunc to provide gzip compression if supported by the client
func GzipHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		h(gw, r)
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
