package organizers

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(pool *pgxpool.Pool) {
	db = pool
}

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

func getSingleOrganizer(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	var o Organizer
	err := db.QueryRow(ctx, "SELECT id, name, org_number FROM organizers WHERE id=$1", id).Scan(&o.Id, &o.Name, &o.OrgNumber)
	if err != nil {
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}
	api.WriteJSON(w, o)
}

func postOrganizer(w http.ResponseWriter, r *http.Request) {
	var o Organizer
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "UPDATE organizers SET name=$1, org_number=$2 WHERE id=$3", o.Name, o.OrgNumber, o.Id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Failed to update organizer", http.StatusInternalServerError)
		return
	}
	api.WriteJSON(w, o)
}

func deleteOrganizer(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.OrganizersPath)
	if isList {
		ctx := r.Context()
		if _, err := db.Exec(ctx, "DELETE FROM organizers"); err != nil {
			http.Error(w, "Failed to delete organizers", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid organizer id", http.StatusBadRequest)
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
