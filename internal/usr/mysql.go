package usr

import (
	"database/sql"
	"errors"
	"fmt"
    "log"

	_ "github.com/go-sql-driver/mysql"
)

type _Mysql struct {
    username string
    password string
    ip string
    port string
    table string
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
