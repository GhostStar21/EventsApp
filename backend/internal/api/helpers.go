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

// WriteJSON sets JSON header and writes v as JSON to the ResponseWriter.
func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
