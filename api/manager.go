package api

import (
	"log"
	"net/http"
	"onward-path/internal/usr"
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
	http.HandleFunc("/loginAdmin", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		if err := xui.Login(username, password); err != nil {
			log.Printf("Calling login api failed: ", err)
		}
	})

	http.HandleFunc("/addClient", func(w http.ResponseWriter, r *http.Request) {
		xui.AddClient(w, r)
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		usr.LoginHandler(w, r)
	})

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		usr.RegisterHandler(w, r)
	})

	http.HandleFunc("/TestUserExistance", func(w http.ResponseWriter, r *http.Request) {
		usr.TestUserExistance(w, r)
	})

	http.HandleFunc("/api/BuyConfig", func(w http.ResponseWriter, r *http.Request) {
		usr.BuyConfigHandler(w, r)
	})

	http.HandleFunc("/api/Auth-check", func(w http.ResponseWriter, r *http.Request) {
		usr.AuthenticateCheckHandler(w, r)
	})

	http.HandleFunc("/api/Get-user-configs", func(w http.ResponseWriter, r *http.Request) {
		usr.GetCurrentConfigHandler(w, r)
	})

	log.Println("API apis has been loaded")
}
