package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseOrganizerRegistrationPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/register-organizer", strings.NewReader(`{"name":"Acme","orgNumber":123456,"type":"Official"}`))
	req.Header.Set("Content-Type", "application/json")

	payload, hasPayload, err := parseOrganizerRegistrationPayload(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !hasPayload {
		t.Fatal("expected a payload to be parsed")
	}
	if payload.Name != "Acme" {
		t.Fatalf("expected organizer name to be parsed, got %q", payload.Name)
	}
	if payload.OrgNumber != 123456 {
		t.Fatalf("expected org number to be parsed, got %d", payload.OrgNumber)
	}
	if payload.Type != "Official" {
		t.Fatalf("expected organizer type to be preserved, got %q", payload.Type)
	}
}

func TestNormalizeOrganizerTypeDefaultsToPersonal(t *testing.T) {
	if got := normalizeOrganizerType(""); got != "Personal" {
		t.Fatalf("expected default type to be Personal, got %q", got)
	}
}
