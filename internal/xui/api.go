package xui

import (
	"fmt"
	"log"
	"onward-path/internal/ipc"
)

var (
	HOST                 = "192.168.109.128"
	PORT                 = 18496
	URI_PATH      string = "t22OMBH6rHZ09Zr/"
	BASE_ENDPOINT        = "panel/api/inbounds/"
)

func Login(username string, password string) error {
	params := map[string]string{
		"username": username,
		"password": password,
	}
	url := fmt.Sprintf("%s:%d/%slogin/", HOST, PORT, URI_PATH)

	result, err := ipc.Post(url, params)

	if err != nil {
		log.Printf("Login of user '%s' failed: '%s'", username, err)
		return err
	}
	log.Printf("Login of user '%s' was successful! | output: '%s'", username, result)

	return nil
}
