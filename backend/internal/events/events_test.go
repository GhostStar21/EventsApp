package events

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"EventsApp/internal/consts"
)

func copyEvents(src []Events) []Events {
	result := make([]Events, len(src))
	copy(result, src)
	return result
}

func restoreEvents(original []Events) {
	events = copyEvents(original)
}

func TestEventsHandler_ListEvents(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	req := httptest.NewRequest(http.MethodGet, "/v1/events/", nil)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got []Events
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != len(original) {
		t.Fatalf("expected %d events, got %d", len(original), len(got))
	}
}

func TestEventsHandler_GetSingleEvent(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	req := httptest.NewRequest(http.MethodGet, "/v1/events/1", nil)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got Events
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Id != 1 || got.Name != "All dress party" {
		t.Fatalf("unexpected event returned: %+v", got)
	}
}

func TestEventsHandler_PostEvents(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	payload := `{"name":"Test Event","isExclusive":true,"event_date":"2026-12-01T00:00:00Z","event_time":"0001-01-01T18:00:00Z","location":"Test","description":"Test event"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", 1)
	ctx = context.WithValue(ctx, "role", string(consts.RoleOrganizer))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var got Events
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Id == 0 {
		t.Fatal("expected non-zero event id")
	}
	if got.Name != "Test Event" {
		t.Fatalf("expected name %q, got %q", "Test Event", got.Name)
	}
	if len(events) != len(original)+1 {
		t.Fatalf("expected events list length %d, got %d", len(original)+1, len(events))
	}
}

func TestFormatEventDateAndTime(t *testing.T) {
	event := Events{
		Date: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		Time: time.Date(1, 1, 1, 18, 30, 0, 0, time.UTC),
	}

	dateValue, timeValue, err := formatEventDateAndTime(event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if dateValue != "2026-06-01" {
		t.Fatalf("expected date %q, got %q", "2026-06-01", dateValue)
	}
	if timeValue != "18:30:00" {
		t.Fatalf("expected time %q, got %q", "18:30:00", timeValue)
	}
}

func TestEventsHandler_PostEvents_ForbiddenForUser(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	payload := `{"name":"Test Event"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/events/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", 1)
	ctx = context.WithValue(ctx, "role", string(consts.RoleUser))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}

func TestEventsHandler_UpdateEvents_AllowedForOrganizer(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	payload := `{"id":1,"name":"Updated Party Name","location":"Oslo"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/events/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "userID", 1)
	ctx = context.WithValue(ctx, "role", string(consts.RoleOrganizer))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestEventsHandler_DeleteEvents_AllowedForOrganizer(t *testing.T) {
	original := copyEvents(events)
	defer restoreEvents(original)

	req := httptest.NewRequest(http.MethodDelete, "/v1/events/1", nil)
	ctx := context.WithValue(req.Context(), "userID", 1)
	ctx = context.WithValue(ctx, "role", string(consts.RoleOrganizer))
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	EventsHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}
