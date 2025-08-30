package xui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"net/http/cookiejar"
)

var (
	Cookie *cookiejar.Jar
	Config *_Config
)

func initCookie() error {
	if Cookie == nil {
		_Cookie, err := cookiejar.New(nil)
		if err != nil {
			log.Println("Failed to create cookie jar: %v", err)
			return err
		}
		Cookie = _Cookie
		log.Println("Cookie has been successfully initilized!")
	}
	return nil
}

func clearCookie() {
	Cookie = nil
}

type XUI struct {
}

func New() *XUI {
	return &XUI{}
}

func (i XUI) Load() error {
	// i.loadAPIs()
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	Cookie = jar
	fmt.Println("Cookie has been initilized")

	Config = newConfig()
	if err = Config.Load(); err != nil {
		fmt.Println("Failed to initilize XUI config: '%v'", err)
		return err
	}
	fmt.Println("XUI config has been initilized")
	var confJson []byte
	if confJson, err = json.Marshal(Config); err != nil {
		log.Printf("Couldn't marshal xui config js: %v", err)
		return err
	}
	fmt.Println(fmt.Printf("Main config IP: '%s'", Config.ServerConfigList[0].Host()))
	fmt.Printf("XUI config: '%s'", string(confJson))

	fmt.Println("XUI module has been loaded")
	return nil
}

func (i XUI) Run() error {
	fmt.Println("XUI module has been run")
	return nil
}

// TODO: del later
func (i XUI) loadAPIs() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		if err := Login(username, password); err != nil {
			log.Printf("Calling login api failed: ", err)
		}
	})

	fmt.Println("XUI apis has been loaded")
}
