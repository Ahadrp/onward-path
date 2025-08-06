package usr

import (
	"log"
	"net/http/cookiejar"
)

type USR struct {
	cookie *cookiejar.Jar
}

func New() *USR {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}

	return &USR{
		cookie: jar,
	}
}

func (u USR) Load() error {
	log.Println("USR module has been loaded")
	return nil
}

func (u USR) Run() error {
	log.Println("USR module has been run")
	return nil
}
