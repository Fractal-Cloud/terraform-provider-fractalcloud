package fractalCloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrganization_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/organizations/org-123" {
			t.Errorf("path = %s, want /organizations/org-123", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Organization{
			Id:          "org-123",
			DisplayName: "My Org",
			Description: "Test organization",
			Status:      "Active",
			AdminsIds:   []string{"admin1"},
			MembersIds:  []string{"member1", "member2"},
		})
	}))
	defer server.Close()

	c := testClient(t, server)
	org, err := c.GetOrganization(context.Background(), "org-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org == nil {
		t.Fatal("organization is nil")
	}
	if org.Id != "org-123" {
		t.Errorf("Id = %q, want %q", org.Id, "org-123")
	}
	if org.DisplayName != "My Org" {
		t.Errorf("DisplayName = %q, want %q", org.DisplayName, "My Org")
	}
	if len(org.MembersIds) != 2 {
		t.Errorf("len(MembersIds) = %d, want 2", len(org.MembersIds))
	}
}

func TestGetOrganization_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	org, err := c.GetOrganization(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if org != nil {
		t.Error("expected nil for 404")
	}
}

func TestGetOrganization_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	c := testClient(t, server)
	_, err := c.GetOrganization(context.Background(), "org-123")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGetOrganization_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	c := testClient(t, server)
	_, err := c.GetOrganization(context.Background(), "org-123")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
