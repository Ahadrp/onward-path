package usr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	USER_TABLE = "user"
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

	if err := login(loginParam); err != nil {
		log.Printf("Login of user '%s' was failed: '%v'", err)
	}
	log.Printf("Login of user '%s' was successful!", loginParam.Email)
}

func Register(w http.ResponseWriter, r *http.Request) {
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

	if err := login(loginParam); err == nil { // Login was successful. So, user is exists!
		log.Printf("User '%s' is already  exists!", loginParam.Email)
		return
	}

	if err := addUser(loginParam); err != nil {
		log.Printf("Couldn't add user '%s' to database: '%v'", loginParam.Email, err)
		return
	}
	log.Printf("Registeration of User '%s' was successful!", loginParam.Email)
}

func addUser(loginParam LoginParam) error {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return fmt.Errorf(err)
	}

	query := fmt.Sprintf("INSERT INTO %s (email, password) VALUES (?, ?)", USER_TABLE)
	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		if _, err := db.Exec(query, loginParam.Email, loginParam.Passwd); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Printf("Couldn't add user '%s' to database: '%v'", loginParam.Email,
			err)
		return err
	}

	return nil
}

func login(loginParam LoginParam) error {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return fmt.Errorf(err)
	}

	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email=? AND password=?", USER_TABLE)
	username := ""
	passwd := ""

	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		err := db.QueryRow(query, loginParam.Email, loginParam.Passwd).Scan(&username, &passwd)
		return err
	}); err != nil {
		log.Printf("No user with username '%s' password '%s'", loginParam.Email,
			loginParam.Passwd)
		return err
	}

	return nil
}

func TestUserExistance(w http.ResponseWriter, r *http.Request) {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return
	}

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

	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email=? AND password=?", USER_TABLE)
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
