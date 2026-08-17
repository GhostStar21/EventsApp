package organizers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func copyOrganizers(src []Organizer) []Organizer {
	result := make([]Organizer, len(src))
	copy(result, src)
	return result
}

func restoreOrganizers(original []Organizer) {
	organizers = copyOrganizers(original)
}

func TestOrganizersHandler_ListOrganizers(t *testing.T) {
	original := copyOrganizers(organizers)
	defer restoreOrganizers(original)

	req := httptest.NewRequest(http.MethodGet, "/v1/organizers/", nil)
	rr := httptest.NewRecorder()

	OrganizersHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got []Organizer
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != len(original) {
		t.Fatalf("expected %d organizers, got %d", len(original), len(got))
	}
}

func TestOrganizersHandler_PostOrganizer(t *testing.T) {
	original := copyOrganizers(organizers)
	defer restoreOrganizers(original)

	payload := `{"name":"New Organizer","orgnumber":999}`
	req := httptest.NewRequest(http.MethodPost, "/v1/organizers/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	OrganizersHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var got Organizer
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Id == 0 {
		t.Fatal("expected non-zero organizer id")
	}
	if got.Name != "New Organizer" {
		t.Fatalf("expected name %q, got %q", "New Organizer", got.Name)
	}
	if len(organizers) != len(original)+1 {
		t.Fatalf("expected organizers list length %d, got %d", len(original)+1, len(organizers))
	}
}
