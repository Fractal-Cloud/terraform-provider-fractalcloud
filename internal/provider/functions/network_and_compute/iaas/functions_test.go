package iaas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// --- VirtualNetwork ---

func TestVirtualNetworkFunction_Metadata(t *testing.T) {
	f := NewVirtualNetworkFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_virtual_network" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_virtual_network", resp.Name)
	}
}

func TestVirtualNetworkFunction_Definition(t *testing.T) {
	f := NewVirtualNetworkFunction()
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

// --- Subnet ---

func TestSubnetFunction_Metadata(t *testing.T) {
	f := NewSubnetFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_subnet" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_subnet", resp.Name)
	}
}

func TestSubnetFunction_Definition(t *testing.T) {
	f := NewSubnetFunction()
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

// --- SecurityGroup ---

func TestSecurityGroupFunction_Metadata(t *testing.T) {
	f := NewSecurityGroupFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_security_group" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_security_group", resp.Name)
	}
}

func TestSecurityGroupFunction_Definition(t *testing.T) {
	f := NewSecurityGroupFunction()
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

// --- VirtualMachine ---

func TestVirtualMachineFunction_Metadata(t *testing.T) {
	f := NewVirtualMachineFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_virtual_machine" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_virtual_machine", resp.Name)
	}
}

func TestVirtualMachineFunction_Definition(t *testing.T) {
	f := NewVirtualMachineFunction()
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

// --- LoadBalancer ---

func TestLoadBalancerFunction_Metadata(t *testing.T) {
	f := NewLoadBalancerFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_load_balancer" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_load_balancer", resp.Name)
	}
}

func TestLoadBalancerFunction_Definition(t *testing.T) {
	f := NewLoadBalancerFunction()
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
