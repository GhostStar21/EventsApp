package events

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(pool *pgxpool.Pool) {
	db = pool
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEvents(w, r)
	case http.MethodPost:
		postEvents(w, r)
	case http.MethodPut:
		updateEvents(w, r)
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
	ctx := r.Context()
	rows, err := db.Query(ctx, "SELECT id, name, is_exclusive, event_date, event_time, location, description FROM events ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to query events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []Events
	for rows.Next() {
		var e Events
		var date time.Time
		var tm time.Time
		if err := rows.Scan(&e.Id, &e.Name, &e.IsExclusive, &date, &tm, &e.Location, &e.Description); err != nil {
			http.Error(w, "Failed to scan event", http.StatusInternalServerError)
			return
		}
		e.Date = date
		e.Time = tm
		out = append(out, e)
	}
	api.WriteJSON(w, out)
}

func getSingleEvent(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	var e Events
	var date time.Time
	var tm time.Time
	err := db.QueryRow(ctx, "SELECT id, name, is_exclusive, event_date, event_time, location, description FROM events WHERE id=$1", id).Scan(&e.Id, &e.Name, &e.IsExclusive, &date, &tm, &e.Location, &e.Description)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	e.Date = date
	e.Time = tm
	api.WriteJSON(w, e)
}

func postEvents(w http.ResponseWriter, r *http.Request) {
	var event Events
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	var id int
	err := db.QueryRow(ctx, "INSERT INTO events (name, is_exclusive, event_date, event_time, location, description) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		event.Name, event.IsExclusive, event.Date, event.Time, event.Location, event.Description).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	event.Id = id
	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, event)
}

func updateEvents(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.EventsPath)
	if isList {
		http.Error(w, "Event id required for update", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}

	var event Events
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if event.Id == 0 {
		event.Id = id
	}
	if event.Id != id {
		http.Error(w, "Event ID mismatch", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "UPDATE events SET name=$1, is_exclusive=$2, event_date=$3, event_time=$4, location=$5, description=$6 WHERE id=$7",
		event.Name, event.IsExclusive, event.Date, event.Time, event.Location, event.Description, event.Id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}
	api.WriteJSON(w, event)
}

func deleteEvents(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.EventsPath)
	if isList {
		// delete all
		ctx := r.Context()
		if _, err := db.Exec(ctx, "DELETE FROM events"); err != nil {
			http.Error(w, "Failed to delete events", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "DELETE FROM events WHERE id=$1", id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
