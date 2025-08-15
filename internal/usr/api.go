package usr

import (
	"log"
	"net/http"
	"bytes"
	"encoding/json"
	"io"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

    var loginParam LoginParam
	bodyBytes, err := io.ReadAll(r.Body)
	if err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&loginParam); err != nil {
		log.Printf("HTTP %d - %s: %s | Error: %v", http.StatusBadRequest, "Invalid JSON body",
			string(bodyBytes), err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

}

func login(loginParam LoginParam) {
    // Mysql.sendQuery()
}
