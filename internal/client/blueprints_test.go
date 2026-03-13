package fractalCloud

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceGroupId_ToString(t *testing.T) {
	id := ResourceGroupId{Type: "Personal", OwnerId: "user1", ShortName: "my-rg"}
	got := id.ToString()
	want := "Personal/user1/my-rg"
	if got != want {
		t.Errorf("ToString() = %q, want %q", got, want)
	}
}

func TestFractalId_ToString(t *testing.T) {
	id := FractalId{
		ResourceGroupId: ResourceGroupId{Type: "Organizational", OwnerId: "org1", ShortName: "rg1"},
		Name:            "my-fractal",
		Version:         "1.0",
	}
	got := id.ToString()
	want := "Organizational/org1/rg1/my-fractal:1.0"
	if got != want {
		t.Errorf("ToString() = %q, want %q", got, want)
	}
}

func testFractalId() FractalId {
	return FractalId{
		ResourceGroupId: ResourceGroupId{Type: "Personal", OwnerId: "user1", ShortName: "rg1"},
		Name:            "test-fractal",
		Version:         "1.0",
	}
}

func TestGetBlueprint_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		wantPath := "/blueprints/Personal/user1/rg1/test-fractal/1.0"
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}

		resp := BlueprintInternal{
			FractalId:   "frac-1",
			IsPrivate:   true,
			Status:      "Active",
			Description: "Test blueprint",
			CreatedAt:   "2026-01-01T00:00:00Z",
			Components: []ComponentInternal{
				{
					Id:                "comp-1",
					Type:              "BigData.PaaS.ComputeCluster",
					DisplayName:       "My Cluster",
					Description:       "A cluster",
					Version:           "v1",
					IsLocked:          true,
					RecreateOnFailure: false,
					Parameters:        map[string]interface{}{"clusterName": "etl", "numWorkers": float64(2)},
					DependenciesIds:   []string{"dep-1", "dep-2"},
					Links: []ComponentLinkInternal{
						{ComponentId: "link-target", Settings: map[string]interface{}{"fromPort": "8080", "protocol": "tcp"}},
					},
					OutputFields: []string{"clusterId", "url"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := testClient(t, server)
	blueprint, err := c.GetBlueprint(context.Background(), testFractalId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blueprint == nil {
		t.Fatal("blueprint is nil")
	}

	if blueprint.FractalId != "frac-1" {
		t.Errorf("FractalId = %q, want %q", blueprint.FractalId, "frac-1")
	}
	if !blueprint.IsPrivate {
		t.Error("IsPrivate = false, want true")
	}
	if blueprint.Description != "Test blueprint" {
		t.Errorf("Description = %q, want %q", blueprint.Description, "Test blueprint")
	}

	if len(blueprint.Components) != 1 {
		t.Fatalf("len(Components) = %d, want 1", len(blueprint.Components))
	}

	comp := blueprint.Components[0]
	if comp.Id != "comp-1" {
		t.Errorf("Component.Id = %q, want %q", comp.Id, "comp-1")
	}
	if comp.Type != "BigData.PaaS.ComputeCluster" {
		t.Errorf("Component.Type = %q, want %q", comp.Type, "BigData.PaaS.ComputeCluster")
	}

	// Pointer fields
	if comp.DisplayName == nil || *comp.DisplayName != "My Cluster" {
		t.Errorf("DisplayName = %v, want 'My Cluster'", comp.DisplayName)
	}
	if comp.Version == nil || *comp.Version != "v1" {
		t.Errorf("Version = %v, want 'v1'", comp.Version)
	}
	if comp.IsLocked == nil || !*comp.IsLocked {
		t.Error("IsLocked should be true")
	}
	if comp.RecreateOnFailure == nil || *comp.RecreateOnFailure {
		t.Error("RecreateOnFailure should be false")
	}

	// Parameters — string values pass through, non-string get JSON-marshaled
	if comp.Parameters["clusterName"] != "etl" {
		t.Errorf("Parameters[clusterName] = %q, want %q", comp.Parameters["clusterName"], "etl")
	}
	if comp.Parameters["numWorkers"] != "2" {
		t.Errorf("Parameters[numWorkers] = %q, want %q", comp.Parameters["numWorkers"], "2")
	}

	// Dependencies
	if len(comp.DependenciesIds) != 2 {
		t.Fatalf("len(DependenciesIds) = %d, want 2", len(comp.DependenciesIds))
	}

	// Links
	if len(comp.Links) != 1 {
		t.Fatalf("len(Links) = %d, want 1", len(comp.Links))
	}
	if comp.Links[0].ComponentId != "link-target" {
		t.Errorf("Link.ComponentId = %q, want %q", comp.Links[0].ComponentId, "link-target")
	}
	if comp.Links[0].Settings["fromPort"] != "8080" {
		t.Errorf("Link.Settings[fromPort] = %q, want %q", comp.Links[0].Settings["fromPort"], "8080")
	}

	// OutputFields
	if len(comp.OutputFields) != 2 {
		t.Fatalf("len(OutputFields) = %d, want 2", len(comp.OutputFields))
	}
}

func TestGetBlueprint_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	blueprint, err := c.GetBlueprint(context.Background(), testFractalId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blueprint != nil {
		t.Error("expected nil blueprint for 404")
	}
}

func TestGetBlueprint_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	c := testClient(t, server)
	_, err := c.GetBlueprint(context.Background(), testFractalId())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGetBlueprint_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	c := testClient(t, server)
	_, err := c.GetBlueprint(context.Background(), testFractalId())
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestCreateBlueprint(t *testing.T) {
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
	err := c.CreateBlueprint(context.Background(), testFractalId(), "desc", true, []Component{
		{Id: "c1", Type: "BigData.PaaS.ComputeCluster"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/blueprints/Personal/user1/rg1/test-fractal/1.0" {
		t.Errorf("path = %s, want /blueprints/Personal/user1/rg1/test-fractal/1.0", gotPath)
	}

	var reqBody CreateBlueprintCommandRequestBody
	if err := json.Unmarshal(gotBody, &reqBody); err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}
	if reqBody.Description != "desc" {
		t.Errorf("Description = %q, want %q", reqBody.Description, "desc")
	}
	if !reqBody.IsPrivate {
		t.Error("IsPrivate = false, want true")
	}
	if len(reqBody.Components) != 1 {
		t.Errorf("len(Components) = %d, want 1", len(reqBody.Components))
	}
}

func TestUpdateBlueprint(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.UpdateBlueprint(context.Background(), testFractalId(), "updated", false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PUT" {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
}

func TestDeleteBlueprint(t *testing.T) {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeleteBlueprint(context.Background(), testFractalId())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

func TestDeleteBlueprint_404IsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	c := testClient(t, server)
	err := c.DeleteBlueprint(context.Background(), testFractalId())
	if err != nil {
		t.Fatalf("delete with 404 should succeed, got error: %v", err)
	}
}
