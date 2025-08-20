package usr

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"onward-path/config"

	"encoding/json"
	"os"
)

type _Mysql struct {
	Username string `json:"username"`
	Password string `json:"passwd"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

func (m *_Mysql) Load() error {
	if err := m.loadConfig(); err != nil {
		log.Printf("Couldn't load Mysql config : %v", err)
		return err
	}
	log.Printf("Mysql has been loaded successfully!")

	return nil
}

func (m *_Mysql) loadConfig() error {
	// Read file
	data, err := os.ReadFile(config.MYSQL_CONFIG)
	if err != nil {
		log.Printf("Couldn't read mysql config file: %v", err)
		return err
	}

	// Parse JSON
	if err := json.Unmarshal(data, m); err != nil {
		log.Printf("Couldn't unmarshal mysql config json: %v", err)
		return err
	}

	return nil
}

func (m *_Mysql) SendQuery(query string, callback func(db *sql.DB) error) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.Username, m.Password, m.Ip, m.Port, m.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Couldn't connect to mysql: %v", err)
		return err
	}
	defer db.Close()

	err = callback(db)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("sth is wrong while sending query to mysql: %v", err)
		}
		return err
	}

	return nil
}
