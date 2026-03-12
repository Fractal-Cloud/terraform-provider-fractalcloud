package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &RelationalDatabaseFunction{}

type RelationalDatabaseFunction struct{}

func NewRelationalDatabaseFunction() function.Function {
	return &RelationalDatabaseFunction{}
}

func (f *RelationalDatabaseFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "relational_database"
}

func (f *RelationalDatabaseFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Relational Database blueprint component",
		Description: "Builds a Relational Database component with the correct type for use in a fractal's components list. " +
			"If dbms_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Relational Database configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"dbms_id":      types.StringType,
					"collation":    types.StringType,
					"charset":      types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type relationalDatabaseConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	DbmsId      types.String `tfsdk:"dbms_id"`
	Collation   types.String `tfsdk:"collation"`
	Charset     types.String `tfsdk:"charset"`
}

func (f *RelationalDatabaseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config relationalDatabaseConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.Collation.IsNull() && !config.Collation.IsUnknown() {
		params["collation"] = config.Collation.ValueString()
	}
	if !config.Charset.IsNull() && !config.Charset.IsUnknown() {
		params["charset"] = config.Charset.ValueString()
	}

	var deps []string
	if !config.DbmsId.IsNull() && !config.DbmsId.IsUnknown() {
		deps = append(deps, config.DbmsId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.RelationalDatabase",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		deps,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
