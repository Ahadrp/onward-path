package api

import (
	"log"
	"net/http"
	"onward-path/internal/xui"
)

type API struct {
}

func New() *API {
	return &API{}
}

func (i API) Load() error {
	i.loadAPIs()

	log.Println("API module has been loaded")
	return nil
}

func (i API) Run() error {
	log.Println("API module has been run")
	return nil
}

func (i API) loadAPIs() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		if err := xui.Login(username, password); err != nil {
			log.Printf("Calling login api failed: ", err)
		}
	})

	http.HandleFunc("/addClient", func(w http.ResponseWriter, r *http.Request) {
		xui.AddClient(w, r)
	})

	log.Println("API apis has been loaded")
}
