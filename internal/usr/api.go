package usr

import (
	"encoding/json"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Handle CORS preflight
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.109.100:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	// Handle CORS preflight
	// w.Header().Set("Access-Control-Allow-Origin", "http://192.168.109.100:5173")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var resp UsrResponseAPI
	token, err := Login(w, r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = "Login has been successful!"
		token = "\"" + token + "\""
		resp.Obj = json.RawMessage(token)
	}

	// Encode struct as JSON and write directly to w
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func BuyConfigHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Handle CORS preflight
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.109.100:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

func AuthenticateCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")
	// Handle CORS preflight
	// w.Header().Set("Access-Control-Allow-Origin", "http://192.168.109.100:5173")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var resp UsrResponseAPI
	_, err := AuthenticateCheck(w, r)
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = "Checking authentication was successful!"
	}

	// Encode struct as JSON and write directly to w
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
