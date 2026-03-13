package fractalCloud

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testPersonalRGId() ResourceGroupId {
	return ResourceGroupId{Type: "Personal", ShortName: "my-rg"}
}

func TestGetPersonalResourceGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/accounts/me/resourcegroups/my-rg" {
			t.Errorf("path = %s, want /accounts/me/resourcegroups/my-rg", r.URL.Path)
		}
		json.NewEncoder(w).Encode(PersonalResourceGroup{
			Id:          ResourceGroupId{Type: "Personal", OwnerId: "user1", ShortName: "my-rg"},
			DisplayName: "My RG",
			Status:      "Active",
			FractalsIds: []string{"f1"},
		})
	}))
	defer server.Close()

	c := testClient(t, server)
	rg, err := c.GetPersonalResourceGroup(context.Background(), testPersonalRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rg == nil {
		t.Fatal("resource group is nil")
	}
	if rg.DisplayName != "My RG" {
		t.Errorf("DisplayName = %q, want %q", rg.DisplayName, "My RG")
	}
	if rg.Id.ShortName != "my-rg" {
		t.Errorf("ShortName = %q, want %q", rg.Id.ShortName, "my-rg")
	}
}

func TestGetPersonalResourceGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	rg, err := c.GetPersonalResourceGroup(context.Background(), testPersonalRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rg != nil {
		t.Error("expected nil for 404")
	}
}

func TestUpsertPersonalResourceGroup(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.UpsertPersonalResourceGroup(context.Background(), PersonalResourceGroup{
		Id:          ResourceGroupId{Type: "Personal", ShortName: "my-rg"},
		DisplayName: "Updated RG",
		Description: "A description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/accounts/me/resourcegroups/my-rg" {
		t.Errorf("path = %s, want /accounts/me/resourcegroups/my-rg", gotPath)
	}

	var reqBody UpsertPersonalResourceGroupRequestBody
	if err := json.Unmarshal(gotBody, &reqBody); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if reqBody.DisplayName != "Updated RG" {
		t.Errorf("DisplayName = %q, want %q", reqBody.DisplayName, "Updated RG")
	}
	if reqBody.Description != "A description" {
		t.Errorf("Description = %q, want %q", reqBody.Description, "A description")
	}
}

func TestDeletePersonalResourceGroup(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeletePersonalResourceGroup(context.Background(), testPersonalRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

func TestDeletePersonalResourceGroup_404IsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeletePersonalResourceGroup(context.Background(), testPersonalRGId())
	if err != nil {
		t.Fatalf("delete with 404 should succeed, got error: %v", err)
	}
}
