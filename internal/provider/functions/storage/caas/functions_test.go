package caas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func TestSearchFunction_Metadata(t *testing.T) {
	f := NewStorageCaasSearchFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_caas_search" {
		t.Errorf("expected name %q, got %q", "storage_caas_search", resp.Name)
	}
}

func TestSearchFunction_Definition(t *testing.T) {
	f := NewStorageCaasSearchFunction()
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

func TestSearchEntityFunction_Metadata(t *testing.T) {
	f := NewStorageCaasSearchEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_caas_search_entity" {
		t.Errorf("expected name %q, got %q", "storage_caas_search_entity", resp.Name)
	}
}

func TestSearchEntityFunction_Definition(t *testing.T) {
	f := NewStorageCaasSearchEntityFunction()
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
