package xui

import (
	"fmt"
	"log"
	"net/http"

	"net/http/cookiejar"
)

var (
	Cookie *cookiejar.Jar
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
