package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error

	connStr := "postgres://artur:0905@localhost:5432/dbweb3?sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to create an object for connecting to database: %v", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	log.Println("Connected to the database successfully")
}
