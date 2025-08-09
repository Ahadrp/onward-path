package usr

import (
	"log"
	"net/http/cookiejar"
)

type USR struct {
	Cookie *cookiejar.Jar
}

func New() *USR {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}

	return &USR{
		Cookie: jar,
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
