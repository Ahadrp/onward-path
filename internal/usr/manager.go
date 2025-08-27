package usr

import (
	"log"
	"net/http/cookiejar"
)

const (
	USER_TABLE    = "user"
	SESSION_TABLE = "sessions"
	SESSION_NAME  = "onwardpath_session_token"
	INBOUND_ID  = 1
)

var (
	Mysql *_Mysql
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
	Mysql = &_Mysql{}
	if err := Mysql.Load(); err != nil {
		log.Println("Couldn't load Mysql: %v", err)
		return err
	}
	log.Println("Mysql has been set-up!")

	log.Println("USR module has been loaded")
	return nil
}

func (u USR) Run() error {
	log.Println("USR module has been run")
	return nil
}
