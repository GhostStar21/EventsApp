package auth

import (
	"context"
	"testing"
)

func TestGetUserIDFromContext_ReturnsUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), "userID", 42)

	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID != 42 {
		t.Fatalf("expected userID 42, got %d", userID)
	}
}

func TestGetUserIDFromContext_MissingUserID(t *testing.T) {
	_, err := GetUserIDFromContext(context.Background())
	if err == nil {
		t.Fatal("expected an error when userID is missing from context")
	}
}
