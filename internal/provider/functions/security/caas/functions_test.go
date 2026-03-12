package caas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// --- ServiceMeshSecurity ---

func TestServiceMeshSecurityFunction_Metadata(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "security_caas_service_mesh_security" {
		t.Errorf("expected name %q, got %q", "security_caas_service_mesh_security", resp.Name)
	}
}

func TestServiceMeshSecurityFunction_Definition(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}
