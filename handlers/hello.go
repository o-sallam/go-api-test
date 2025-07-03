package handlers

import (
	"encoding/json"
	"net/http"

	"go-api-test/models"
)

// HelloWorldHandler handles the /hello endpoint
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Message: "Hello World!",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(response)
}
