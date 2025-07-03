package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-api-test/handlers"
	"go-api-test/utils"
)

var staticRoot = "wwwroot"

func main() {
	// Load HTML at startup (no minify)
	htmlBytes, err := os.ReadFile(staticRoot + "/index.html")
	if err != nil {
		log.Fatalf("Failed to load index.html: %v", err)
	}
	handlers.SetPortfolioHTML(string(htmlBytes))

	// Load CSS at startup (no minify)
	cssBytes, err := os.ReadFile(staticRoot + "/style.css")
	if err != nil {
		log.Fatalf("Failed to load style.css: %v", err)
	}
	cssContent := string(cssBytes)

	mux := http.NewServeMux()

	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(staticRoot+"/img"))))
	mux.HandleFunc("/majallat-althaqafa.html", handlers.HomeHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/majallat-althaqafa.html", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/hello", handlers.HelloWorldHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/favicon.ico", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, staticRoot+"/favicon.ico")
	}))
	mux.HandleFunc("/style.css", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write([]byte(cssContent))
	}))
	mux.HandleFunc("/robots.txt", utils.GzipHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, r, staticRoot+"/robots.txt")
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Try: http://localhost:%s/hello\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("Portfolio page: http://localhost:%s/portfolio\n", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
