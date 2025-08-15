package usr

import (
	"database/sql"
	"errors"
	"fmt"
    "log"

	_ "github.com/go-sql-driver/mysql"
    "internal/onward-path/config"

    "encoding/json"
	"os"
)

type _Mysql struct {
    username string `json:"username"`
    password string `json:"passwd"`
    ip string `json:"ip"`
    port string `json:"port"`
    table string `json:"table"`
}

func (m _Mysql) LoadConfig() error {
    // Read file
	data, err := os.ReadFile(config.MYSQL_CONFIG)
	if err != nil {
		log.Printf("Couldn't read mysql config file: %v", err)
		return err
	}

	// Parse JSON
	if err := json.Unmarshal(data, &m); err != nil {
		log.Printf("Couldn't unmarshal mysql config json: %v", err)
		return err
	}

    return nil
}

func (m _Mysql) SendQuery(query string, callback func(db *sql.DB) error) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.username, m.password, m.ip, m.port, m.table)
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
