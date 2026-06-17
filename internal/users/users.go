package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"log"
	"net/http"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		postUsers(w, r)
	case http.MethodPut:
		updateUsers(w, r)
	case http.MethodDelete:
		deleteUsers(w, r)
	default:
		http.Error(w, "This method ( "+r.Method+" ) is not supported", http.StatusNotImplemented)
		return
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		listUsers(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	getSingleUser(w, r, id)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	api.WriteJSON(w, users)
}

func getSingleUser(w http.ResponseWriter, r *http.Request, id int) {
	for _, user := range users {
		if user.Id == id {
			api.WriteJSON(w, user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Println("Received POST Request")

	if u.Id == 0 {
		max := 0
		for _, ex := range users {
			if ex.Id > max {
				max = ex.Id
			}
		}
		u.Id = max + 1
	}
	users = append(users, u)
	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, u)
}

func updateUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		http.Error(w, "User id required for update", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if u.Id == 0 {
		u.Id = id
	}
	if u.Id != id {
		http.Error(w, "User ID mismatch", http.StatusBadRequest)
		return
	}

	for i := range users {
		if users[i].Id == id {
			users[i] = u
			api.WriteJSON(w, u)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		users = []User{}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	for i, u := range users {
		if u.Id == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}
