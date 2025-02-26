package main

import (
	"gin-app/database"
	"gin-app/routes"
	"log"
)

func main() {
	// Initialize the database
	database.InitDB()
	defer database.GetDB().Close()

	// Set up the Gin router using the routes package
	r := routes.SetupRouter()

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
