package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySql driver
)

// InitDB - Database connection
func InitDB() *sql.DB {
	dbUser := os.Getenv("MYSQL_USER")
	dbPassw := os.Getenv("MYSQL_PASSWORD")
	dbDatabase := os.Getenv("MYSQL_DATABASE")
	connection := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s", dbUser, dbPassw, dbDatabase)
	db, err := sql.Open("mysql", connection)

	if err != nil {
		log.Printf(connection)
		panic(err.Error())
	}
	return db
}
