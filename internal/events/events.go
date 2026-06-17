package events

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"log"
	"net/http"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEvents(w, r)
	case http.MethodPost:
		postEvents(w, r)
	case http.MethodDelete:
		deleteEvents(w, r)
	default:
		http.Error(w, "This method ( "+r.Method+" ) is not supported", http.StatusNotImplemented)
		return
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.EventsPath)
	if isList {
		listEvents(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}
	getSingleEvent(w, r, id)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	api.WriteJSON(w, events)
}

func getSingleEvent(w http.ResponseWriter, r *http.Request, id int) {
	for _, event := range events {
		if event.Id == id {
			api.WriteJSON(w, event)
			return
		}
	}
	http.Error(w, "Event not found", http.StatusNotFound)
}

func postEvents(w http.ResponseWriter, r *http.Request) {
	var event Events

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Println("Received POST Request")

	// assign ID if not provided
	if event.Id == 0 {
		max := 0
		for _, e := range events {
			if e.Id > max {
				max = e.Id
			}
		}
		event.Id = max + 1
	}

	events = append(events, event)

	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, event)
}

func deleteEvents(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.EventsPath)
	if isList {
		events = []Events{}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}

	for i, event := range events {
		if event.Id == id {
			events = append(events[:i], events[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Event not found", http.StatusNotFound)
}
