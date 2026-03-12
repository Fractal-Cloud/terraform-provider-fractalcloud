package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// --- APIGateway ---

func TestAPIGatewayFunction_Metadata(t *testing.T) {
	f := NewPaaSAPIGatewayFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "api_management_paas_api_gateway" {
		t.Errorf("expected name %q, got %q", "api_management_paas_api_gateway", resp.Name)
	}
}

func TestAPIGatewayFunction_Definition(t *testing.T) {
	f := NewPaaSAPIGatewayFunction()
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
