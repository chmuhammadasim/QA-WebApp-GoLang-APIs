package main

import (
	"log"
	"net/http"
	"os"
	"qa-app/db"
	"qa-app/routes"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	// Connect to the database
	db.Connect()

	// Initialize the router
	router := routes.InitRoutes()

	// Add logging to the router
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Log server start before starting
	log.Println("Server starting on port 8080")

	// Start the server
	err = http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
