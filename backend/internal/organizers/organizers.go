package organizers

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Sets the global database connection pool.
func SetDB(pool *pgxpool.Pool) {
	db = pool
}

// Handles different HTTP methods.
func OrganizersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrganizer(w, r)
	case http.MethodPost:
		postOrganizer(w, r)
	case http.MethodPut:
		updateOrganizer(w, r)
	case http.MethodDelete:
		deleteOrganizer(w, r)
	default:
		http.Error(w, "This method ( "+r.Method+" ) is not supported", http.StatusNotImplemented)
		return
	}
}

// Deletes an organizer.
func deleteOrganizer(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.OrganizersPath)
	if isList {
		deleteAllOrganizers(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid organizer id", http.StatusBadRequest)
		return
	}
	deleteSingleOrganizer(w, r, id)
}

// Deletes all organizers (only available for admin).
func deleteAllOrganizers(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		organizers = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ctx := r.Context()
	if _, err := db.Exec(ctx, "DELETE FROM organizers"); err != nil {
		http.Error(w, "Failed to delete organizers", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Deletes a single organizer.
func deleteSingleOrganizer(w http.ResponseWriter, r *http.Request, id int) {
	if db == nil {
		for i, o := range organizers {
			if o.Id == id {
				organizers = append(organizers[:i], organizers[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "DELETE FROM organizers WHERE id=$1", id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Lists either all organizers or a single organizer.
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

// Lists all organizers that exist in the database.
func listOrganizers(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		api.WriteJSON(w, organizers)
		return
	}
	ctx := r.Context()
	rows, err := db.Query(ctx, "SELECT id, name, org_number FROM organizers ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to query organizers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []Organizer
	for rows.Next() {
		var o Organizer
		if err := rows.Scan(&o.Id, &o.Name, &o.OrgNumber); err != nil {
			http.Error(w, "Failed to scan organizer", http.StatusInternalServerError)
			return
		}
		out = append(out, o)
	}
	api.WriteJSON(w, out)
}

// Get a single organizer from the id
func getSingleOrganizer(w http.ResponseWriter, r *http.Request, id int) {
	if db == nil {
		for _, o := range organizers {
			if o.Id == id {
				api.WriteJSON(w, o)
				return
			}
		}
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	ctx := r.Context()
	var o Organizer
	err := db.QueryRow(ctx, "SELECT id, name, org_number FROM organizers WHERE id=$1", id).Scan(&o.Id, &o.Name, &o.OrgNumber)
	if err != nil {
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	api.WriteJSON(w, o)
}

// Register an organizer into the database.
func postOrganizer(w http.ResponseWriter, r *http.Request) {
	var o Organizer
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if db == nil {
		o.Id = len(organizers) + 1
		organizers = append(organizers, o)
		w.WriteHeader(http.StatusCreated)
		api.WriteJSON(w, o)
		return
	}
	ctx := r.Context()
	var id int
	err := db.QueryRow(ctx, "INSERT INTO organizers (name, org_number) VALUES ($1,$2) RETURNING id", o.Name, o.OrgNumber).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create organizer", http.StatusInternalServerError)
		return
	}
	o.Id = id
	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, o)
}

// Update the event / make changes in the organizer data.
func updateOrganizer(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.OrganizersPath)
	if isList {
		http.Error(w, "Organizer id required for update", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, "Invalid organizer id", http.StatusBadRequest)
		return
	}

	var o Organizer
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if o.Id == 0 {
		o.Id = id
	}

	if o.Id != id {
		http.Error(w, "Organizer ID mismatch", http.StatusBadRequest)
		return
	}
	if db == nil {
		for i, existing := range organizers {
			if existing.Id == id {
				organizers[i] = o
				api.WriteJSON(w, o)
				return
			}
		}
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "UPDATE organizers SET name=$1, org_number=$2 WHERE id=$3", o.Name, o.OrgNumber, o.Id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Failed to update organizer", http.StatusInternalServerError)
		return
	}
	api.WriteJSON(w, o)
}
