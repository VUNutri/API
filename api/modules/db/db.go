package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	log.Printf("Trying to connect")
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/nutri")
	if err != nil {
		panic(err.Error())
	}
	return db
}
