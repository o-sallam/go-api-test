package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-api-test/handlers"
	"go-api-test/utils"
)

func main() {
	// Load and minify HTML at startup
	htmlBytes, err := os.ReadFile("wwwroot/index.html")
	if err != nil {
		log.Fatalf("Failed to load index.html: %v", err)
	}
	minified := utils.MinifyHTML(string(htmlBytes))
	handlers.SetPortfolioHTML(minified)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.PortfolioHandler)
	mux.HandleFunc("/hello", handlers.HelloWorldHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, "wwwroot/favicon.ico")
	})
	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, "wwwroot/style.css")
	})

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
