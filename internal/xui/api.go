package xui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"onward-path/internal/ipc"
)

var (
	HOST                  = "192.168.109.128"
	PORT                  = 18496
	URI_PATH       string = "t22OMBH6rHZ09Zr/"
	BASE_ENDPOINT         = "panel/api/inbounds/"
	ADMIN_USERNAME        = "root"
	ADMIN_PASSWD          = "123"
)

func Login(username string, password string) error {
	if err := initCookie(); err != nil {
		log.Println("Login failed because: '%v'", err)
		return err
	}

	params := map[string]string{
		"username": username,
		"password": password,
	}
	url := fmt.Sprintf("%s:%d/%slogin/", HOST, PORT, URI_PATH)

	result, err := ipc.PostLogin(url, params, Cookie)

	if err != nil {
		log.Printf("Login of user '%s' failed: '%s'", username, err)
		clearCookie()
		return err
	}
	log.Printf("Login of user '%s' was successful! | output: '%s'", username, result)

	return nil
}

func AddClient(w http.ResponseWriter, r *http.Request) {
	if err := Login(ADMIN_USERNAME, ADMIN_PASSWD); err != nil {
		log.Printf("Login of user '%s' failed: '%s'", ADMIN_USERNAME, err)
		return
	}

	addClient(w, r)
}

func GetClient(email string) (json.RawMessage, error) {
	if err := Login(ADMIN_USERNAME, ADMIN_PASSWD); err != nil {
		log.Printf("Login of user '%s' failed: '%s'", ADMIN_USERNAME, err)
		return json.RawMessage{}, err
	}

	var _client json.RawMessage
	var err error
	if _client, err = getClient(email); err != nil {
		log.Printf("Couldn't get user '%s': %v", email, err)
		return json.RawMessage{}, err
	}

	return _client, nil
}

func getClient(email string) (json.RawMessage, error) {
	endPoint := "getClientTraffics/"
	url := fmt.Sprintf("%s:%d/%s%s%s", HOST, PORT, URI_PATH, BASE_ENDPOINT, endPoint)

	result, err := ipc.Get(url, email, Cookie)
	if err != nil {
		log.Printf("Sending Get request failed: '%v'", err)
		return json.RawMessage{}, err
	}
	log.Printf(result)

	var xuiResp XUIResponse
	err = json.Unmarshal([]byte(result), &xuiResp)
	if err != nil {
		log.Printf("Couldn't parse xui response: '%v'", err)
		return json.RawMessage{}, err
	}

	return xuiResp.Obj, nil
}

func addClient(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s:%d/%s%saddClient/", HOST, PORT, URI_PATH, BASE_ENDPOINT)
	// find user base on session. assume we've found it.
	// TODO: check if user exist with this email.
	if r.Method != http.MethodPost {
		errTxt := "Method Not Allowed"
		log.Printf("HTTP %d - %s", http.StatusMethodNotAllowed, errTxt)
		http.Error(w, errTxt, http.StatusMethodNotAllowed)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var addClientRequestExternalAPI AddClientRequestExternalAPI
	bodyBytes, err := io.ReadAll(r.Body)
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&addClientRequestExternalAPI); err != nil {
		log.Printf("HTTP %d - %s: %s", http.StatusBadRequest, "Invalid JSON body",
			string(bodyBytes))
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	addClientRequestExternalAPI.Settings.Clients[0].ID = uuid.New().String()

	internalClientJson, err := json.Marshal(addClientRequestExternalAPI.Settings)
	if err != nil {
		log.Printf("json client error")
	}

	addClientRequestInternalAPI := AddClientRequestInternalAPI{
		ID:       addClientRequestExternalAPI.ID,
		Settings: string(internalClientJson),
	}

	/*
	   jsonClient, err := json.Marshal(addClientRequest)
	   if err != nil {
	       log.Printf("Failed to convert client to json: ", err)
	       return
	   }
	*/

	criaJson, err := json.Marshal(addClientRequestInternalAPI)
	if err != nil {
		log.Printf("json error")
		return
	}

	result, err := ipc.Post(url, string(criaJson), Cookie)
	if err != nil {
		log.Printf("Failed to convert client to json: ", err)
		return
	}
	log.Printf("Client '%s' was added successfully! | output: '%s'", addClientRequestExternalAPI.Settings.Clients[0].Email, result)

}
