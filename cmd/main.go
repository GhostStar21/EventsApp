package main

import (
	"EventsApp/internal/auth"
	"EventsApp/internal/consts"
	"EventsApp/internal/events"
	"EventsApp/internal/organizers"
	"EventsApp/internal/postgres"
	"EventsApp/internal/users"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Main function that handles endpoints and initiates/connects to database and port.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	} else {
		log.Println("$PORT has been set: ", port)
	}

	// Connect to database
	ctx := context.Background()
	db, err := postgres.NewPool(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Provide DB pool to packages that need it
	events.SetDB(db)
	users.SetDB(db)
	organizers.SetDB(db)

	defer db.Close()

	log.Println("Successfully connected to database")

	// Initiate router
	router := http.NewServeMux()

	// Handle the various endpoints
	router.HandleFunc(consts.RegisterPath, auth.Register(db))
	router.HandleFunc(consts.LoginPath, auth.Login(db))
	router.HandleFunc(consts.MePath, auth.AuthMiddleware(users.MeHandler))
	router.HandleFunc(consts.EventsPath, events.EventsHandler)
	router.HandleFunc(consts.UsersPath, auth.AuthMiddleware(users.UsersHandler))
	router.HandleFunc(consts.OrganizersPath, organizers.OrganizersHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
