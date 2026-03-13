package fractalCloud

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testOrgRGId() ResourceGroupId {
	return ResourceGroupId{Type: "Organizational", OwnerId: "org-1", ShortName: "my-rg"}
}

func TestGetOrganizationalResourceGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/organizations/org-1/resourcegroups/my-rg" {
			t.Errorf("path = %s, want /organizations/org-1/resourcegroups/my-rg", r.URL.Path)
		}
		json.NewEncoder(w).Encode(OrganizationalResourceGroup{
			Id:          ResourceGroupId{Type: "Organizational", OwnerId: "org-1", ShortName: "my-rg"},
			DisplayName: "Org RG",
			Status:      "Active",
			MembersIds:  []string{"m1", "m2"},
			TeamsIds:    []string{"t1"},
			ManagersIds: []string{"mgr1"},
		})
	}))
	defer server.Close()

	c := testClient(t, server)
	rg, err := c.GetOrganizationalResourceGroup(context.Background(), testOrgRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rg == nil {
		t.Fatal("resource group is nil")
	}
	if rg.DisplayName != "Org RG" {
		t.Errorf("DisplayName = %q, want %q", rg.DisplayName, "Org RG")
	}
	if len(rg.MembersIds) != 2 {
		t.Errorf("len(MembersIds) = %d, want 2", len(rg.MembersIds))
	}
	if len(rg.TeamsIds) != 1 {
		t.Errorf("len(TeamsIds) = %d, want 1", len(rg.TeamsIds))
	}
}

func TestGetOrganizationalResourceGroup_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	rg, err := c.GetOrganizationalResourceGroup(context.Background(), testOrgRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rg != nil {
		t.Error("expected nil for 404")
	}
}

func TestUpsertOrganizationalResourceGroup(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(202)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.UpsertOrganizationalResourceGroup(context.Background(), OrganizationalResourceGroup{
		Id:          ResourceGroupId{Type: "Organizational", OwnerId: "org-1", ShortName: "my-rg"},
		DisplayName: "Updated Org RG",
		Description: "Org description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/organizations/org-1/resourcegroups/my-rg" {
		t.Errorf("path = %s, want /organizations/org-1/resourcegroups/my-rg", gotPath)
	}

	var reqBody UpsertOrganizationalResourceGroupRequestBody
	if err := json.Unmarshal(gotBody, &reqBody); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if reqBody.DisplayName != "Updated Org RG" {
		t.Errorf("DisplayName = %q, want %q", reqBody.DisplayName, "Updated Org RG")
	}
}

func TestDeleteOrganizationalResourceGroup(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeleteOrganizationalResourceGroup(context.Background(), testOrgRGId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

func TestDeleteOrganizationalResourceGroup_404IsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeleteOrganizationalResourceGroup(context.Background(), testOrgRGId())
	if err != nil {
		t.Fatalf("delete with 404 should succeed, got error: %v", err)
	}
}
