package routes

import (
	"go-api-test/handlers"
	"go-api-test/utils"
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, staticRoot string, cssContent string) {
	// --- Static file endpoints ---
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(staticRoot+"/img"))))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.FileServer(http.Dir(staticRoot+"/fonts")).ServeHTTP(w, r)
	})))
	mux.HandleFunc("/favicon.ico", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, staticRoot+"/favicon.ico")
	}))
	mux.Handle("/css/", http.StripPrefix("/css/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.FileServer(http.Dir(staticRoot+"/css")).ServeHTTP(w, r)
	})))
	mux.HandleFunc("/robots.txt", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, r, staticRoot+"/robots.txt")
	}))
	mux.HandleFunc("/google4fe8d22092105d8e.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, staticRoot+"/google4fe8d22092105d8e.html")
	})
	mux.Handle("/js/", http.StripPrefix("/js/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.FileServer(http.Dir(staticRoot+"/js")).ServeHTTP(w, r)
	})))

	// --- Page endpoints ---
	mux.HandleFunc("/hello", handlers.HelloWorldHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/post-partial-html/", handlers.PostPartialHTMLHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handlers.HomeHandler(w, r)
			return
		}
		// Redirect /slug/ to /slug (remove trailing slash)
		if strings.HasSuffix(r.URL.Path, "/") {
			clean := strings.TrimSuffix(r.URL.Path, "/")
			http.Redirect(w, r, clean, http.StatusMovedPermanently)
			return
		}
		handlers.PostHandler(w, r)
	})
}
