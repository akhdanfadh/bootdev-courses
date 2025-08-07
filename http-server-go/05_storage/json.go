package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse is a struct to hold error messages in JSON format
type errorResponse struct {
	Error string `json:"error"`
}

// respondJSON is a utility function to respond with JSON data given payload
func respondJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload) // encode the payload
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(500) // since encoding is server-side
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
