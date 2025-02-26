package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// Initialize the database connection
func InitDB() {
	var err error
	connStr := "user=postgres dbname=todo_db sslmode=disable password=postgress host=localhost port=5432"

	// connStr := "user=postgres dbname=todo_db sslmode=disable password=yourpassword host=postgres port=5432"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to the database!")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}
