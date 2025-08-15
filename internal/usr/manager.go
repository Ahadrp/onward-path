package usr

import (
	"log"
	"net/http/cookiejar"
    "internal/onward-path/config"
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
    if err := Mysql.LoadConfig(); err != nil {
        log.Println("Couldn't load Mysql config: %v", err)
        return err
    }
    log.Println("Mysql config has been set-up!")

	log.Println("USR module has been loaded")
	return nil
}

func (u USR) Run() error {
	log.Println("USR module has been run")
	return nil
}
