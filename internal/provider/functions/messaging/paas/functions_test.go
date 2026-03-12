package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func TestBrokerFunction_Metadata(t *testing.T) {
	f := NewMessagingPaasBrokerFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_paas_broker" {
		t.Errorf("expected name %q, got %q", "messaging_paas_broker", resp.Name)
	}
}

func TestBrokerFunction_Definition(t *testing.T) {
	f := NewMessagingPaasBrokerFunction()
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

func TestEntityFunction_Metadata(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_paas_entity" {
		t.Errorf("expected name %q, got %q", "messaging_paas_entity", resp.Name)
	}
}

func TestEntityFunction_Definition(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
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
