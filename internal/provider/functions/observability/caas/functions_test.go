package caas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// --- Monitoring ---

func TestMonitoringFunction_Metadata(t *testing.T) {
	f := NewCaaSMonitoringFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_monitoring" {
		t.Errorf("expected name %q, got %q", "observability_caas_monitoring", resp.Name)
	}
}

func TestMonitoringFunction_Definition(t *testing.T) {
	f := NewCaaSMonitoringFunction()
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

// --- Tracing ---

func TestTracingFunction_Metadata(t *testing.T) {
	f := NewCaaSTracingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_tracing" {
		t.Errorf("expected name %q, got %q", "observability_caas_tracing", resp.Name)
	}
}

func TestTracingFunction_Definition(t *testing.T) {
	f := NewCaaSTracingFunction()
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

// --- Logging ---

func TestLoggingFunction_Metadata(t *testing.T) {
	f := NewCaaSLoggingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_logging" {
		t.Errorf("expected name %q, got %q", "observability_caas_logging", resp.Name)
	}
}

func TestLoggingFunction_Definition(t *testing.T) {
	f := NewCaaSLoggingFunction()
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
