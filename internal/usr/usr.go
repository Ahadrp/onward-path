package usr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"onward-path/internal/xui"
)

type User struct {
	xui.Client
}

func Register(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		// http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return "", fmt.Errorf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
	}

	// Set response header
	// w.Header().Set("Content-Type", "application/json")

	var loginParam LoginParam
	bodyBytes, err := io.ReadAll(r.Body)
	if err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&loginParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("HTTP %d - %s: %s | Error: %v", http.StatusBadRequest, "Invalid JSON body",
			string(bodyBytes), err)
		// http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return "", fmt.Errorf("Invalid JSON body")
	}
	defer r.Body.Close()

	if err := findUser(loginParam); err == nil { // Login was successful. So, user is exists!
		errText := fmt.Sprintf("User '%s' is already  exists!", loginParam.Email)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	if err := addUser(loginParam); err != nil {
		errText := fmt.Sprintf("Couldn't add user '%s' to database: '%v'", loginParam.Email, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}
	log.Printf("Registeration of User '%s' was successful!", loginParam.Email)

	return "", nil
}

func Login(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		// http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return "", fmt.Errorf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
	}

	// Set response header
	// w.Header().Set("Content-Type", "application/json")

	var loginParam LoginParam
	bodyBytes, err := io.ReadAll(r.Body)
	if err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&loginParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("HTTP %d - %s: %s | Error: %v", http.StatusBadRequest, "Invalid JSON body",
			string(bodyBytes), err)
		// http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return "", fmt.Errorf("Invalid JSON body")
	}
	defer r.Body.Close()

	if err := findUser(loginParam); err != nil {
		errText := fmt.Sprintf("Can not find user '%s': '%v'", loginParam.Email, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	var token string
	if token, err = login(loginParam); err != nil {
		errText := fmt.Sprintf("Login of user '%s' was failed: '%v'", err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	// return token to browser as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_NAME,
		Value:    token,
		Path:     "/",
		HttpOnly: true, // JS cannot read it
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // required with None
	})

	log.Printf("Login of user '%s' was successful!", loginParam.Email)

	return token, nil
}

func BuyConfig(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		// http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return "", fmt.Errorf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// No valid token...
		w.WriteHeader(http.StatusUnauthorized)
		errText := fmt.Sprintf("Missing session token")
		// http.Error(w, "Missing session token", http.StatusUnauthorized)
		log.Printf(errText)
		return "", fmt.Errorf(errText)
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Set response header
	// w.Header().Set("Content-Type", "application/json")

	// var addClientRequestExternalAPI AddClientRequestExternalAPI
	var addClientParam AddClientParam
	bodyBytes, err := io.ReadAll(r.Body)
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&addClientParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("HTTP %d - %s: %s", http.StatusBadRequest, "Invalid JSON body",
			string(bodyBytes))
		// http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return "", fmt.Errorf("Invalid JSON body")
	}
	defer r.Body.Close()

	email, err := findUserByToken(token)
	if err != nil {
		errText := fmt.Sprintf("Couldn't find any user with token '%s': ", token, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}
	addClientParam.Email = email

	var _client json.RawMessage
	if _client, err = xui.GetClient(email, addClientParam.Server); err != nil {
		errText := fmt.Sprintf("Get client '%s' from server '%d' failed: '%v'", email, addClientParam.Server, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	var client xui.GetClientResponse
	err = json.Unmarshal(_client, &client)
	if err != nil {
		errText := fmt.Sprintf("Failed to process client '%s' json: '%v'", email, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	if client.Email == "" {
		if err := buyConfig(&addClientParam); err != nil {
			errText := fmt.Sprintf("Failed to buy account for user '%s': %v", email, err)
			log.Println(errText)
			return "", fmt.Errorf(errText)
		}
		log.Printf("Account for '%s' has been created!", email)
		return "", nil
	} else {
		errText := fmt.Sprintf("User '%s' has already an account!", client.Email)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}
}

func GetCurrentConfig(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		// http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return "", fmt.Errorf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// No valid token...
		w.WriteHeader(http.StatusUnauthorized)
		errText := fmt.Sprintf("Missing session token")
		// http.Error(w, "Missing session token", http.StatusUnauthorized)
		log.Printf(errText)
		return "", fmt.Errorf(errText)
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Set response header
	// w.Header().Set("Content-Type", "application/json")

	email, err := findUserByToken(token)
	if err != nil {
		errText := fmt.Sprintf("Couldn't find any user with token '%s': ", token, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	var clientConfigsRaw json.RawMessage
	if clientConfigsRaw, err = xui.GetUserConfigs(email); err != nil {
		errText := fmt.Sprintf("Get configs of client '%s' failed: '%v'", email, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	var clientConfigList xui.CurrentConfigList
	err = json.Unmarshal(clientConfigsRaw, &clientConfigList)
	if err != nil {
		errText := fmt.Sprintf("Failed to process configs of client '%s' json: '%v'", email, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	if len(clientConfigList.CurrentConfigs) <= 0 {
		errText := fmt.Sprintf("User '%s' has no config!", email)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	} else {
		log.Printf("Configs of user '%s' has been found successfully!", email)
		return string([]byte(clientConfigsRaw)), nil
	}
}

func AuthenticateCheck(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		// http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return "", fmt.Errorf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// No valid token...
		w.WriteHeader(http.StatusUnauthorized)
		errText := fmt.Sprintf("Missing session token")
		// http.Error(w, "Missing session token", http.StatusUnauthorized)
		log.Printf(errText)
		return "", fmt.Errorf(errText)
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	exist, err := checkSessionExistance(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errText := fmt.Sprintf("Error in finding token '%s': '%v'", token, err)
		log.Println(errText)
		return "", fmt.Errorf(errText)
	}

	if exist {
		w.WriteHeader(http.StatusOK)
		log.Printf("Token '%s' is valid!", token)
		return "", nil
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		errText := fmt.Sprintf("No suck a token: '%s'", token)
		log.Printf(errText)
		return "", fmt.Errorf(errText)
	}
}

func addUser(loginParam LoginParam) error {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return fmt.Errorf(err)
	}

	hashedPasswd := SHA256(loginParam.Passwd)

	query := fmt.Sprintf("INSERT INTO %s (email, password) VALUES (?, ?)", USER_TABLE)
	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		if _, err := db.Exec(query, loginParam.Email, hashedPasswd); err != nil {
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

func findUser(loginParam LoginParam) error {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return fmt.Errorf(err)
	}

	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email=? AND password=?", USER_TABLE)
	email := ""
	passwd := ""

	hashedPasswd := SHA256(loginParam.Passwd)

	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		err := db.QueryRow(query, loginParam.Email, hashedPasswd).Scan(&email, &passwd)
		return err
	}); err != nil {
		log.Printf("No user with email '%s' password '%s'", loginParam.Email,
			loginParam.Passwd)
		return err
	}

	return nil
}

func findUserByToken(token string) (string, error) {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return "", fmt.Errorf(err)
	}

	// query to db for finding user
	query := fmt.Sprintf("SELECT email FROM %s WHERE token=?", SESSION_TABLE)
	email := ""

	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		err := db.QueryRow(query, token).Scan(&email)
		return err
	}); err != nil {
		log.Printf("Couldn't send query to db: '%v'", err)
		return "", err
	}

	return email, nil
}

func login(loginParam LoginParam) (string, error) {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return "", fmt.Errorf(err)
	}

	// create token
	token := GenerateRandomToken()

	// add token to db
	query := fmt.Sprintf("INSERT INTO %s (email, token) VALUES (?, ?)", SESSION_TABLE)
	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		_, err := db.Exec(query, loginParam.Email, token)
		if err != nil {
			log.Printf("Couldn't add session of user '%s' to database: '%v'", loginParam.Email, err)
		}
		return err
	}); err != nil {
		log.Printf("Couldn't send query to db: '%v'", loginParam.Email, err)
		return "", err
	}

	return token, nil
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
	email := ""
	passwd := ""

	if err = Mysql.SendQuery(query, func(db *sql.DB) error {
		err := db.QueryRow(query, loginParam.Email, loginParam.Passwd).Scan(&email, &passwd)
		return err
	}); err != nil {
		log.Printf("No user with email '%s' password '%s'", loginParam.Email,
			loginParam.Passwd)
		return
	}
	log.Printf("User with email %s exist!", loginParam.Email)
}

func buyConfig(addClientParam *AddClientParam) error {
	xuiAddClientParam := xui.AddClientRequestExternalAPI{
		Server: addClientParam.Server,
		ID:     INBOUND_ID,
		Settings: xui.SettingsDecoded{
			Clients: []xui.ClientParam{
				{
					Email:      addClientParam.Email,
					Flow:       addClientParam.Flow,
					TotalGB:    addClientParam.Total,
					ExpiryTime: addClientParam.ExpiryTime,
				},
			},
		},
	}

	if err := xui.AddClientInternal(xuiAddClientParam); err != nil {
		return err
	}

	return nil
}

func checkSessionExistance(token string) (bool, error) {
	if Mysql == nil {
		err := "Mysql obj hasn't been initilized!"
		log.Printf(err)
		return false, fmt.Errorf(err)
	}

	// query to db for finding user
	query := fmt.Sprintf("SELECT token FROM %s WHERE token=?", SESSION_TABLE)
	_token := ""

	if err := Mysql.SendQuery(query, func(db *sql.DB) error {
		err := db.QueryRow(query, token).Scan(&_token)
		return err
	}); err != nil {
		if err == sql.ErrNoRows {
			// Not found
			log.Printf("Token not found: '%s'", token)
			return false, nil
		}
		// Some other DB error
		log.Printf("Couldn't send query to db: '%v'", err)
		return false, err
	}

	return true, nil
}
