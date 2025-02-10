package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
    "os"
    "github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get database credentials from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    //dbName := os.Getenv("DB_NAME")

    // Define a flag to accept the database name from the command line
    dbName := flag.String("dbname", "", "Name of the database to create")
    flag.Parse()

    // Check if the database name is provided
    if *dbName == "" {
        log.Fatal("Please provide a database in your")
    }

    // Construct the MySQL connection string
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbHost, dbPort)

    // Open a connection to the MySQL server
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }

     // Create the database
     query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", *dbName)
     _, err = db.Exec(query)
     if err != nil {
         log.Fatalf("Failed to create database: %v", err)
     }
 
     fmt.Printf("Database '%s' created successfully!\n", *dbName)

    useDBQuery := fmt.Sprintf("USE %s", *dbName)
	_, err = db.Exec(useDBQuery)
	if err != nil {
		log.Fatalf("Failed to switch to database '%s': %v", *dbName, err)
	}

	// Create the 'todo' table if it does not exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS todo (
			id INT AUTO_INCREMENT PRIMARY KEY,
			task VARCHAR(255) NOT NULL,
			status ENUM('pending', 'completed') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create 'todo' table: %v", err)
	}
	fmt.Println("Table 'todo' created successfully!")
    defer db.Close()
}