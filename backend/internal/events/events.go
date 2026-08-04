package events

import (
	"EventsApp/internal/api"
	"EventsApp/internal/auth"
	"EventsApp/internal/consts"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Sets the global database connection pool.
func SetDB(pool *pgxpool.Pool) {
	db = pool
}

// Handles different HTTP methods.
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

// Deletes an event.
func deleteEvents(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.EventsPath)
	if isList {
		deleteAllEvents(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid event id", http.StatusBadRequest)
		return
	}
	deleteSingleEvent(w, r, id)

}

// Deletes all events (only available for admin).
func deleteAllEvents(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value("role")
	role, _ := roleVal.(string)

	if role != string(consts.RoleAdmin) {
		http.Error(w, "Only admins can delete all events", http.StatusForbidden)
		return
	}

	if db == nil {
		events = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ctx := r.Context()
	if _, err := db.Exec(ctx, "DELETE FROM events"); err != nil {
		http.Error(w, "Failed to delete events", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Deletes a single event.
func deleteSingleEvent(w http.ResponseWriter, r *http.Request, id int) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	roleVal := r.Context().Value("role")
	role, _ := roleVal.(string)

	if role != string(consts.RoleOrganizer) && role != string(consts.RoleAdmin) {
		http.Error(w, "Only organizers can delete events", http.StatusForbidden)
		return
	}

	if db == nil {
		for i, e := range events {
			if e.Id == id {
				events = append(events[:i], events[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	ctx := r.Context()

	var exists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE id=$1)", id).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if role != string(consts.RoleAdmin) {
		var isOwner bool
		err = db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM event_organizers eo
				JOIN organizer_member om ON eo.organizer_id = om.organizer_id
				WHERE eo.event_id = $1 AND om.user_id = $2
			)
		`, id, userID).Scan(&isOwner)

		if err != nil || !isOwner {
			http.Error(w, "Forbidden: You can only delete events created by your organizer", http.StatusForbidden)
			return
		}
	}

	_, err = db.Exec(ctx, "DELETE FROM event_organizers WHERE event_id=$1", id)
	if err != nil {
		http.Error(w, "Failed to remove event organizer links", http.StatusInternalServerError)
		return
	}

	cmdTag, err := db.Exec(ctx, "DELETE FROM events WHERE id=$1", id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Lists either all events or a single event.
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

// Lists all events that exist in the database.
func listEvents(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		api.WriteJSON(w, events)
		return
	}

	ctx := r.Context()
	rows, err := db.Query(ctx, `
		SELECT e.id, e.name, e.is_exclusive, e.event_date, e.event_time, e.location, e.description, e.is_registration,
		       (SELECT eo.organizer_id FROM event_organizers eo WHERE eo.event_id = e.id ORDER BY eo.organizer_id LIMIT 1) AS organizer_id,
		       (SELECT o.name FROM event_organizers eo JOIN organizers o ON eo.organizer_id = o.id WHERE eo.event_id = e.id ORDER BY eo.organizer_id LIMIT 1) AS organizer_name
		FROM events e
		ORDER BY e.id
	`)
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
		var orgID sql.NullInt32
		var orgName sql.NullString

		if err := rows.Scan(&e.Id, &e.Name, &e.IsExclusive, &date, &tm, &e.Location, &e.Description, &e.IsRegistration, &orgID, &orgName); err != nil {
			http.Error(w, "Failed to scan event", http.StatusInternalServerError)
			return
		}
		e.Date = date
		e.Time = tm
		if orgID.Valid {
			v := int(orgID.Int32)
			e.OrganizerId = &v
		}
		if orgName.Valid {
			e.OrganizerName = orgName.String
		}

		out = append(out, e)
	}
	api.WriteJSON(w, out)
}

// Get a single event from the id
func getSingleEvent(w http.ResponseWriter, r *http.Request, id int) {
	if db == nil {
		for _, e := range events {
			if e.Id == id {
				api.WriteJSON(w, e)
				return
			}
		}
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	ctx := r.Context()
	var e Events
	var date time.Time
	var tm time.Time
	var orgID sql.NullInt32
	var orgName sql.NullString

	err := db.QueryRow(ctx, `
		SELECT e.id, e.name, e.is_exclusive, e.event_date, e.event_time, e.location, e.description, e.is_registration,
		       (SELECT eo.organizer_id FROM event_organizers eo WHERE eo.event_id = e.id ORDER BY eo.organizer_id LIMIT 1) AS organizer_id,
		       (SELECT o.name FROM event_organizers eo JOIN organizers o ON eo.organizer_id = o.id WHERE eo.event_id = e.id ORDER BY eo.organizer_id LIMIT 1) AS organizer_name
		FROM events e
		WHERE e.id = $1
	`, id).Scan(&e.Id, &e.Name, &e.IsExclusive, &date, &tm, &e.Location, &e.Description, &e.IsRegistration, &orgID, &orgName)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	e.Date = date
	e.Time = tm
	if orgID.Valid {
		v := int(orgID.Int32)
		e.OrganizerId = &v
	}
	if orgName.Valid {
		e.OrganizerName = orgName.String
	}
	api.WriteJSON(w, e)
}

func formatEventDateAndTime(event Events) (string, string, error) {
	if event.Date.IsZero() {
		return "", "", fmt.Errorf("event date is required")
	}

	dateValue := event.Date.Format("2006-01-02")
	var timeValue string
	if event.Time.IsZero() {
		timeValue = "00:00:00"
	} else {
		timeValue = event.Time.Format("15:04:05")
	}

	return dateValue, timeValue, nil
}

// Register an event into the database.
func postEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	roleVal := r.Context().Value("role")
	role, _ := roleVal.(string)

	if role != string(consts.RoleOrganizer) && role != string(consts.RoleAdmin) {
		http.Error(w, "Only organizers can create events", http.StatusForbidden)
		return
	}

	var event Events
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if db == nil {
		event.Id = len(events) + 1
		events = append(events, event)
		w.WriteHeader(http.StatusCreated)
		api.WriteJSON(w, event)
		return
	}

	ctx := r.Context()

	var targetOrgID int
	if event.OrganizerId != nil && *event.OrganizerId > 0 {
		targetOrgID = *event.OrganizerId
		if role != string(consts.RoleAdmin) {
			var exists bool
			err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM organizer_member WHERE user_id=$1 AND organizer_id=$2)", userID, targetOrgID).Scan(&exists)
			if err != nil || !exists {
				http.Error(w, "You are not a member of this organizer", http.StatusForbidden)
				return
			}
		}
	} else {
		err := db.QueryRow(ctx, "SELECT organizer_id FROM organizer_member WHERE user_id=$1 ORDER BY organizer_id LIMIT 1", userID).Scan(&targetOrgID)
		if err != nil && role != string(consts.RoleAdmin) {
			http.Error(w, "No organizer associated with your account", http.StatusBadRequest)
			return
		}
	}

	dateValue, timeValue, err := formatEventDateAndTime(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow(ctx, "INSERT INTO events (name, is_exclusive, event_date, event_time, location, description, is_registration) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		event.Name, event.IsExclusive, dateValue, timeValue, event.Location, event.Description, event.IsRegistration).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	event.Id = id

	if targetOrgID > 0 {
		_, _ = db.Exec(ctx, "INSERT INTO event_organizers (event_id, organizer_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", id, targetOrgID)
		event.OrganizerId = &targetOrgID

		var orgName string
		_ = db.QueryRow(ctx, "SELECT name FROM organizers WHERE id=$1", targetOrgID).Scan(&orgName)
		event.OrganizerName = orgName
	}

	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, event)
}

// Update the event / make changes in the event.
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

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	roleVal := r.Context().Value("role")
	role, _ := roleVal.(string)

	if role != string(consts.RoleOrganizer) && role != string(consts.RoleAdmin) {
		http.Error(w, "Only organizers can update events", http.StatusForbidden)
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

	if db == nil {
		for i, e := range events {
			if e.Id == id {
				events[i] = event
				api.WriteJSON(w, event)
				return
			}
		}
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	ctx := r.Context()

	var exists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE id=$1)", id).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if role != string(consts.RoleAdmin) {
		var isOwner bool
		err = db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM event_organizers eo
				JOIN organizer_member om ON eo.organizer_id = om.organizer_id
				WHERE eo.event_id = $1 AND om.user_id = $2
			)
		`, id, userID).Scan(&isOwner)

		if err != nil || !isOwner {
			http.Error(w, "Forbidden: You can only update events created by your organizer", http.StatusForbidden)
			return
		}
	}

	cmdTag, err := db.Exec(ctx, "UPDATE events SET name=$1, is_exclusive=$2, event_date=$3, event_time=$4, location=$5, description=$6, is_registration=$7 WHERE id=$8",
		event.Name, event.IsExclusive, event.Date, event.Time, event.Location, event.Description, event.IsRegistration, event.Id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	getSingleEvent(w, r, id)
}
