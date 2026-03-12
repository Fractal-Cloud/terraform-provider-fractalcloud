package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func TestContainerPlatformFunction_Metadata(t *testing.T) {
	f := NewContainerPlatformFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_paas_container_platform" {
		t.Errorf("expected name %q, got %q", "network_and_compute_paas_container_platform", resp.Name)
	}
}

func TestContainerPlatformFunction_Definition(t *testing.T) {
	f := NewContainerPlatformFunction()
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
