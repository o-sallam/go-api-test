package handlers

import (
	"encoding/json"
	"net/http"

	"manage-system.api/models"
)

// HealthHandler handles the /health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Message: "API is running",
		Status:  "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
