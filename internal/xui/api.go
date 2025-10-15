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
	"strconv"
	"strings"
)

// TODO: no need, rm later.
var (
	HOST                  = "192.168.109.128"
	PORT                  = 18496
	URI_PATH       string = "t22OMBH6rHZ09Zr/"
	BASE_ENDPOINT         = "panel/api/inbounds/"
	ADMIN_USERNAME        = "root"
	ADMIN_PASSWD          = "123"
)

// TODO: rm it later. No need.
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

func LoginWithServerID(serverID int) error {
	if len(Config.ServerConfigList) == 0 {
		log.Printf("No server has been defined")
		return fmt.Errorf("Sorry! No server is available now!")
	}

	for _, serverConfig := range Config.ServerConfigList {
		if serverConfig.id == strconv.Itoa(serverID) {
			if err := initCookie(); err != nil {
				log.Println("Login failed because: '%v'", err)
				return err
			}

			params := map[string]string{
				"username": serverConfig.adminUser,
				"password": serverConfig.adminPass,
			}
			url := fmt.Sprintf(
				"%s:%d/%slogin/",
				serverConfig.host,
				serverConfig.port,
				serverConfig.uriPath)

			result, err := ipc.PostLogin(url, params, Cookie)

			if err != nil {
				log.Printf("Login of user '%s' failed: '%s'", serverConfig.adminUser, err)
				clearCookie()
				return err
			}
			log.Printf("Login of user '%s' was successful! | output: '%s'", serverConfig.adminUser, result)
			return nil
		}
	}

	return fmt.Errorf("No such server: '%d'", serverID)
}

// TODO: rm later
func AddClient(w http.ResponseWriter, r *http.Request) {
	if err := Login(ADMIN_USERNAME, ADMIN_PASSWD); err != nil {
		log.Printf("Login of user '%s' failed: '%s'", ADMIN_USERNAME, err)
		return
	}

	addClient(w, r)
}

func AddClientInternal(addClientRequestExternalAPI AddClientRequestExternalAPI) error {
	if err := LoginWithServerID(addClientRequestExternalAPI.Server); err != nil {
		log.Printf("Login of admin failed: '%v'", err)
		return fmt.Errorf("Admin login failed")
	}

	if err := addClient_Internal(addClientRequestExternalAPI); err != nil {
		return err
	}

	return nil
}

// TODO: rm it later.
func GetClient(email string, serverID int) (json.RawMessage, bool, error) {
	if err := LoginWithServerID(serverID); err != nil {
		log.Printf("Login of admin failed: '%v'", err)
		return json.RawMessage{}, false, fmt.Errorf("Admin login failed")
	}

	serverIDString := strconv.Itoa(serverID) // int -> string
	for _, server := range Config.ServerConfigList {
		if server.id != serverIDString {
			// This is not the requested server.
			continue
		}

		var client GetClientResponse
		var hasUser bool
		var err error
		if hasUser, err = getClientWithServerConfig(email, &server, &client); err != nil {
			log.Printf("Couldn't get client with email '%s' from server '%s': '%v'",
				email, server.host, err)
			return json.RawMessage{}, false, err
		}

		if !hasUser {
			log.Printf("No user with email '%s' in server '%s'",
				email, server.host)
			return json.RawMessage{}, false, nil
		}

		clientByte, err := json.Marshal(client)
		if err != nil {
			log.Printf("Couldn't convert user of '%s' to byte: '%v'",
				email, err)
			return json.RawMessage{}, true, err
		}

		return json.RawMessage(clientByte), true, nil
	}

	errText := fmt.Sprintf("No server with id '%d'", serverID)
	return json.RawMessage{}, false, fmt.Errorf(errText)
}

func GetUserConfigs(email string) (json.RawMessage, error) {
	var currentConfigList CurrentConfigList
	for _, server := range Config.ServerConfigList {
		if err := getClientWithServer(email, &server, &currentConfigList); err != nil {
			log.Printf("Couldn't get user of '%s' from server '%s': '%v'", email, server.host, err)
		}
	}

	currentConfigListByte, err := json.Marshal(currentConfigList)
	if err != nil {
		log.Printf("Couldn't convert user of '%s' to byte: '%v'",
			email, err)
		return json.RawMessage{}, err
	}

	return json.RawMessage(currentConfigListByte), nil
}

// TODO: rm it later.
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

// TODO: rm later.
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

func addClient_Internal(addClientRequestExternalAPI AddClientRequestExternalAPI) error {
	url, err := createURL(addClientRequestExternalAPI.Server)
	if err != nil {
		log.Printf("Failed to create URL: '%v'", err)
		return err
	}

	addClientRequestExternalAPI.Settings.Clients[0].ID = uuid.New().String()

	internalClientJson, err := json.Marshal(addClientRequestExternalAPI.Settings)
	if err != nil {
		log.Printf("json client error")
		return err
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
		return err
	}

	result, err := ipc.Post(url, string(criaJson), Cookie)
	if err != nil {
		log.Printf("Failed to convert client to json: '%v'", err)
		return err
	}

	var xuiResp XUIResponse
	err = json.Unmarshal([]byte(result), &xuiResp)
	if err != nil {
		log.Printf("Couldn't parse xui response: '%v'", err)
		return err
	}

	if xuiResp.Success {
		log.Printf("Client '%s' was added successfully! | output: '%s'", addClientRequestExternalAPI.Settings.Clients[0].Email, result)
		return nil
	} else {
		log.Printf("Failed to add client '%s' | output: '%s'", addClientRequestExternalAPI.Settings.Clients[0].Email, result)
		return fmt.Errorf(xuiResp.Message)
	}
}

func createURL(serverID int) (string, error) {
	if len(Config.ServerConfigList) == 0 {
		log.Printf("No server has been defined")
		return "", fmt.Errorf("Sorry! No server is available now!")
	}

	for _, serverConfig := range Config.ServerConfigList {
		if serverConfig.id == strconv.Itoa(serverID) {
			return fmt.Sprintf(
				"%s:%d/%s%saddClient/",
				serverConfig.host,
				serverConfig.port,
				serverConfig.uriPath,
				serverConfig.baseEndpoint), nil
		}
	}

	return "", fmt.Errorf("No server with id '%d'!", serverID)
}

func getClientWithServer(email string, serverConf *serverConfig, currentConfigList *CurrentConfigList) error {

	serverID, err := strconv.Atoi(serverConf.id) // string -> int
	if err != nil {
		fmt.Println("convert error:", err)
		return err
	}

	if err := LoginWithServerID(serverID); err != nil {
		log.Printf("Admin login to server '%s' was failed: '%v'", err)
		return err
	}

	var client GetClientResponse
	var hasUser bool
	if hasUser, err = getClientWithServerConfig(email, serverConf, &client); err != nil {
		log.Printf("Couldn't get client with email '%s' from server '%s': '%v'",
			email, serverConf.host, err)
		return err
	}

	if !hasUser {
		log.Printf("No user with email '%s' in server '%s'",
			email, serverConf.host)
		return err
	}

	var inbound Inbound
	if err = getInbound(serverConf, &inbound); err != nil {
		log.Printf("Couldn't get inbound of server '%s': '%v'", serverConf.host, err)
		return err
	}

	for _, cli := range inbound.Settings.Clients { // set uuid. To create URL
		if cli.Email == client.Client.Email {
			client.Client.UUID = cli.ID
			break
		}
	}

	currentConfig := CurrentConfig{
		Inbound:      inbound,
		ClientConfig: client.Client,
	}

	currentConfigList.CurrentConfigs = append(currentConfigList.CurrentConfigs, currentConfig)

	return nil
}

func getClientWithServerConfig(email string, serverConf *serverConfig, clientConfig *GetClientResponse) (bool, error) {
	endPoint := "getClientTraffics/"
	url := fmt.Sprintf("%s:%d/%s%s%s", serverConf.host, serverConf.port,
		serverConf.uriPath, serverConf.baseEndpoint, endPoint)

	result, err := ipc.Get(url, email, Cookie)
	if err != nil {
		log.Printf("Sending Get request failed: '%v'", err)
		return false, err
	}
	log.Printf(result)

	var xuiResp XUIResponse
	err = json.Unmarshal([]byte(result), &xuiResp)
	if err != nil {
		log.Printf("Couldn't parse xui response: '%v'", err)
		return false, err
	}

	// TODO: bug?
	if string([]byte(xuiResp.Obj)) == "null" { // no client with the email in this server
		errText := fmt.Sprintf("No User with email '%s' in server '%s'", email, serverConf.host)
		log.Println(errText)
		return false, nil
	}

	err = json.Unmarshal([]byte(xuiResp.Obj), &clientConfig)
	if err != nil {
		log.Printf("Couldn't parse get client response: '%v'", err)
		return true, err
	}

	return true, nil
}

func getInbound(serverConf *serverConfig, inbound *Inbound) error {
	endPoint := "get/"
	url := fmt.Sprintf("%s:%d/%s%s%s", serverConf.host, serverConf.port,
		serverConf.uriPath, serverConf.baseEndpoint, endPoint)

	result, err := ipc.Get(url, serverConf.id, Cookie)
	if err != nil {
		log.Printf("Sending Get request failed: '%v'", err)
		return err
	}
	log.Printf(result)

	var xuiResp XUIResponse
	err = json.Unmarshal([]byte(result), &xuiResp)
	if err != nil {
		log.Printf("Couldn't parse xui response: '%v'", err)
		return err
	}

	if string([]byte(xuiResp.Obj)) == "null" { // no client with the email in this server
		errText := fmt.Sprintf("No inbound  with id '%s' in server '%s'", serverConf.id, serverConf.host)
		log.Println(errText)
		return fmt.Errorf(errText)
	}

	resultByte := []byte(xuiResp.Obj)
	cleanResult := strings.ReplaceAll(string(resultByte), "\\n", "")
	cleanResult = strings.ReplaceAll(cleanResult, "\\", "")
	cleanResult = strings.ReplaceAll(cleanResult, "\"{", "{")
	cleanResult = strings.ReplaceAll(cleanResult, "}\"", "}")

	log.Printf(cleanResult)

	err = json.Unmarshal([]byte(cleanResult), inbound)
	if err != nil {
		log.Printf("Couldn't parse get inbound response: '%v'", err)
		return err
	}
	inbound.IP = serverConf.host

	return nil
}
