package saas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// --- Unmanaged ---

func TestUnmanagedFunction_Metadata(t *testing.T) {
	f := NewSaaSUnmanagedFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "security_saas_unmanaged" {
		t.Errorf("expected name %q, got %q", "security_saas_unmanaged", resp.Name)
	}
}

func TestUnmanagedFunction_Definition(t *testing.T) {
	f := NewSaaSUnmanagedFunction()
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
