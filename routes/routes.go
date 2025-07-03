package routes

import (
	"go-api-test/handlers"
	"go-api-test/utils"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, staticRoot string, cssContent string) {
	// --- Static file endpoints ---
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(staticRoot+"/img"))))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(staticRoot+"/fonts"))))
	mux.HandleFunc("/favicon.ico", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, staticRoot+"/favicon.ico")
	}))
	mux.HandleFunc("/style.css", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write([]byte(cssContent))
	}))
	mux.HandleFunc("/header.css", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, staticRoot+"/header.css")
	}))
	mux.HandleFunc("/footer.css", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, staticRoot+"/footer.css")
	}))
	mux.HandleFunc("/robots.txt", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, r, staticRoot+"/robots.txt")
	}))

	// --- Page endpoints ---
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/hello", handlers.HelloWorldHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
}
