package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func TestDistributedDataProcessingFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_distributed_data_processing" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_distributed_data_processing", resp.Name)
	}
}

func TestDistributedDataProcessingFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
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

func TestComputeClusterFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_compute_cluster" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_compute_cluster", resp.Name)
	}
}

func TestComputeClusterFunction_Definition(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
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

func TestDataProcessingJobFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_data_processing_job" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_data_processing_job", resp.Name)
	}
}

func TestDataProcessingJobFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
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

func TestMlExperimentFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasMlExperimentFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_ml_experiment" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_ml_experiment", resp.Name)
	}
}

func TestMlExperimentFunction_Definition(t *testing.T) {
	f := NewBigdataPaasMlExperimentFunction()
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

func TestDatalakeFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDatalakeFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_datalake" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_datalake", resp.Name)
	}
}

func TestDatalakeFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDatalakeFunction()
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
