package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-api-test/handlers"
	"go-api-test/routes"
)

const staticRoot = "wwwroot"

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
	routes.RegisterRoutes(mux, staticRoot, cssContent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Try: http://localhost:%s/\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
