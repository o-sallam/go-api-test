package utils

import (
	"compress/gzip"
	"net/http"
	"os"
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

// Show404 serves the custom 404 page with long cache headers
func Show404(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	page, err := os.ReadFile("views/404.html")
	if err != nil {
		w.Write([]byte("404 - Not Found"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}
