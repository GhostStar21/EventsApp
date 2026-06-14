package events

import (
	"net/http"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEvents(w, r)
	case http.MethodPost:
		postEvents(w, r)
	default:
		http.Error(w, "This method ( " + r.Method + " ) is not supported", http.StatusNotImplemented)
		return
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	
}

func postEvents(w http.ResponseWriter, r *http.Request) {

}

