package xui

import (
	"fmt"
	"log"
	"net/http"
)

type XUI struct {
}

func New() *XUI {
	return &XUI{}
}

func (i XUI) Load() error {
	// i.loadAPIs()

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
