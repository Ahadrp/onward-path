package usr

import (
	"encoding/json"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var resp UsrResponseAPI
	_, err := Register(w, r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = "Registeration has been successful!"
	}

	// Encode struct as JSON and write directly to w
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var resp UsrResponseAPI
	_, err := Login(w, r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = "Login has been successful!"
	}

	// Encode struct as JSON and write directly to w
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func BuyConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var resp UsrResponseAPI
	_, err := BuyConfig(w, r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = "Login has been successful!"
	}

	// Encode struct as JSON and write directly to w
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
