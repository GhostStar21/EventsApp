package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ExtractIDFromRequest returns (id, isList, err).
// isList is true when the request targets the collection (no ID path).
func ExtractIDFromRequest(r *http.Request, basePath string) (int, bool, error) {
	path := strings.TrimPrefix(r.URL.Path, basePath)
	path = strings.Trim(path, "/")
	if path == "" {
		return 0, true, nil
	}
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, false, err
	}
	return id, false, nil
}

// Gives access control to the frontend such that it can perform
// actions on the backend.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// WriteJSON sets JSON header and writes v as JSON to the ResponseWriter.
func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
