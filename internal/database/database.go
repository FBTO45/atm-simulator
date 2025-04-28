package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	// Update with your MySQL credentials
	DB, err = sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/atm_simulator?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL database")
}
