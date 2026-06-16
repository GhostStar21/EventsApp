package main

import (
	"EventsApp/internal/consts"
	"EventsApp/internal/events"
	"log"
	"net/http"
	"os"

)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	} else {
		log.Println("$PORT has been set: ", port)
	}

	router := http.NewServeMux()

	router.HandleFunc(consts.EventsPath, events.EventsHandler)

	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Printf("Failed to start server: %v", err)
	}

}
