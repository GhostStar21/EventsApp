package main

import (
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

	router := http.NewServeMux()

	router.HandleFunc(consts.EventsPath, events.EventsHandler)
	router.HandleFunc(consts.UsersPath, users.UsersHandler)
	router.HandleFunc(consts.OrganizersPath, organizers.OrganizersHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
