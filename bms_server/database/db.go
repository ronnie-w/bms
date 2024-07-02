package database

import (
	"database/sql"
	"log"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

const (
	DRIVER_NAME = "sqlite3"
	DATA_SRC    = "file:bms.db?_pragma_key=232DEB3C8E74B426F7FAA44DAC52E47EBC7FD8EAEADE188676D3764FC1A3A624&_pragma_cipher_page_size=4096"
	//DATA_SRC    = "bms.db"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open(DRIVER_NAME, DATA_SRC)
	if err != nil {
		log.Fatalln(err)
	}
}

func Conn() *sql.DB {
	return db
}

func Exists(table, uniqueKey, uniqueVal string) bool {
	var existCount int
	row := db.QueryRow(`SELECT COUNT(*) AS existCount FROM ` + table + ` WHERE ` + uniqueKey + `="` + uniqueVal + `"`)
	row.Scan(&existCount)

	if existCount > 0 {
		return true
	}

	return false
}
