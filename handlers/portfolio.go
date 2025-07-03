package handlers

import (
	"net/http"
)

var minifiedPortfolioHTML string

// SetPortfolioHTML sets the minified HTML in memory
func SetPortfolioHTML(html string) {
	minifiedPortfolioHTML = html
}

// PortfolioHandler serves the minified HTML from memory
func PortfolioHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(minifiedPortfolioHTML))
}
