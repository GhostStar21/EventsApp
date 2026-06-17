package organizers

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"log"
	"net/http"
)

func OrganizersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrganizer(w, r)
	case http.MethodPost:
		postOrganizer(w, r)
	case http.MethodDelete:
		deleteOrganizer(w, r)
	default:
		http.Error(w, "This method ( "+r.Method+" ) is not supported", http.StatusNotImplemented)
		return
	}
}

func getOrganizer(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.OrganizersPath)
	if isList {
		listOrganizers(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid organizer id", http.StatusBadRequest)
		return
	}
	getSingleOrganizer(w, r, id)
}

func listOrganizers(w http.ResponseWriter, r *http.Request) {
	api.WriteJSON(w, organizers)
}

func getSingleOrganizer(w http.ResponseWriter, r *http.Request, id int) {
	for _, organizer := range organizers {
		if organizer.Id == id {
			api.WriteJSON(w, organizer)
			return
		}
	}
	http.Error(w, "Organizer not found", http.StatusNotFound)
}
func postOrganizer(w http.ResponseWriter, r *http.Request) {
	var o Organizer
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Println("Received POST Request")

	if o.Id == 0 {
		max := 0
		for _, ex := range organizers {
			if ex.Id > max {
				max = ex.Id
			}
		}
		o.Id = max + 1
	}
	organizers = append(organizers, o)
	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, o)
}

func deleteOrganizer(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.OrganizersPath)
	if isList {
		organizers = []Organizer{}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid organizer id", http.StatusBadRequest)
		return
	}
	for i, o := range organizers {
		if o.Id == id {
			organizers = append(organizers[:i], organizers[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Organizer not found", http.StatusNotFound)
}
