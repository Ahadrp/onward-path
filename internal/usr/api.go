package usr

import (
	"log"
    "fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"io"
	"database/sql"
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

func TestUserExistance(w http.ResponseWriter, r *http.Request) {
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

    query := fmt.Sprintf("SELECT email, password FROM user WHERE email=? AND password=?")
    username := ""
    passwd := ""

    if err = Mysql.SendQuery(query, func(db *sql.DB) error {
        err := db.QueryRow(query, loginParam.Email, loginParam.Passwd).Scan(&username, &passwd)
        return err
        }); err != nil {
        log.Printf("No user with username '%s' password '%s'", loginParam.Email,
        loginParam.Passwd)
        return
    }
    log.Printf("User with email %s exist!", loginParam.Email)
}

