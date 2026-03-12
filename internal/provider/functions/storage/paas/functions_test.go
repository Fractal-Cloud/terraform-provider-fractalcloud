package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func TestFilesAndBlobsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasFilesAndBlobsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_files_and_blobs" {
		t.Errorf("expected name %q, got %q", "storage_paas_files_and_blobs", resp.Name)
	}
}

func TestFilesAndBlobsFunction_Definition(t *testing.T) {
	f := NewStoragePaasFilesAndBlobsFunction()
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

func TestRelationalDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasRelationalDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_relational_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_relational_dbms", resp.Name)
	}
}

func TestRelationalDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasRelationalDbmsFunction()
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

func TestRelationalDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasRelationalDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_relational_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_relational_database", resp.Name)
	}
}

func TestRelationalDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasRelationalDatabaseFunction()
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

func TestDocumentDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasDocumentDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_document_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_document_dbms", resp.Name)
	}
}

func TestDocumentDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasDocumentDbmsFunction()
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

func TestDocumentDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasDocumentDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_document_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_document_database", resp.Name)
	}
}

func TestDocumentDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasDocumentDatabaseFunction()
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

func TestColumnOrientedDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasColumnOrientedDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_column_oriented_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_column_oriented_dbms", resp.Name)
	}
}

func TestColumnOrientedDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasColumnOrientedDbmsFunction()
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

func TestColumnOrientedEntityFunction_Metadata(t *testing.T) {
	f := NewStoragePaasColumnOrientedEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_column_oriented_entity" {
		t.Errorf("expected name %q, got %q", "storage_paas_column_oriented_entity", resp.Name)
	}
}

func TestColumnOrientedEntityFunction_Definition(t *testing.T) {
	f := NewStoragePaasColumnOrientedEntityFunction()
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

func TestKeyValueDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasKeyValueDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_key_value_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_key_value_dbms", resp.Name)
	}
}

func TestKeyValueDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasKeyValueDbmsFunction()
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

func TestKeyValueEntityFunction_Metadata(t *testing.T) {
	f := NewStoragePaasKeyValueEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_key_value_entity" {
		t.Errorf("expected name %q, got %q", "storage_paas_key_value_entity", resp.Name)
	}
}

func TestKeyValueEntityFunction_Definition(t *testing.T) {
	f := NewStoragePaasKeyValueEntityFunction()
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

func TestGraphDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasGraphDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_graph_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_graph_dbms", resp.Name)
	}
}

func TestGraphDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasGraphDbmsFunction()
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

func TestGraphDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasGraphDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_graph_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_graph_database", resp.Name)
	}
}

func TestGraphDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasGraphDatabaseFunction()
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
