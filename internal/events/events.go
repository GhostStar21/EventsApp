package events

import (
	"encoding/json"
	"EventsApp/internal/consts"
	"net/http"
	"strconv"
	"strings"
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
		http.Error(w, "This method ( " + r.Method + " ) is not supported", http.StatusNotImplemented)
		return
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, consts.EventsPath)
	path = strings.Trim(path, "/")
	if path == "" {
		json.NewEncoder(w).Encode(events)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}

	for _, event := range events {
		if event.Id == id {
			json.NewEncoder(w).Encode(event)
			return
		}
	}

	http.Error(w, "Event not found", http.StatusNotFound)
}

func postEvents(w http.ResponseWriter, r *http.Request) {

}

func deleteEvents(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, consts.EventsPath)
	path = strings.Trim(path, "/")
	if path == "" {
		events = []Events{}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := strconv.Atoi(path)
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

