package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &DocumentDatabaseFunction{}

type DocumentDatabaseFunction struct{}

func NewDocumentDatabaseFunction() function.Function {
	return &DocumentDatabaseFunction{}
}

func (f *DocumentDatabaseFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "document_database"
}

func (f *DocumentDatabaseFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Document Database blueprint component",
		Description: "Builds a Document Database component with the correct type for use in a fractal's components list. If dbms_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Document Database configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"dbms_id":      types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type documentDatabaseConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	DbmsId      types.String `tfsdk:"dbms_id"`
}

func (f *DocumentDatabaseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config documentDatabaseConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	if !config.DbmsId.IsNull() && !config.DbmsId.IsUnknown() {
		deps = append(deps, config.DbmsId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.DocumentDatabase",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		nil,
		deps,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
