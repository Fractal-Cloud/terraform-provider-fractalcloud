package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &RelationalDbmsFunction{}

type RelationalDbmsFunction struct{}

func NewRelationalDbmsFunction() function.Function {
	return &RelationalDbmsFunction{}
}

func (f *RelationalDbmsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "relational_dbms"
}

func (f *RelationalDbmsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Relational DBMS Platform blueprint component",
		Description: "Builds a Relational DBMS Platform component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Relational DBMS configuration",
				AttributeTypes: map[string]attr.Type{
					"id":             types.StringType,
					"display_name":   types.StringType,
					"description":    types.StringType,
					"version":        types.StringType,
					"engine_version": types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type relationalDbmsConfig struct {
	Id            types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
	Version       types.String `tfsdk:"version"`
	EngineVersion types.String `tfsdk:"engine_version"`
}

func (f *RelationalDbmsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config relationalDbmsConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.EngineVersion.IsNull() && !config.EngineVersion.IsUnknown() {
		params["version"] = config.EngineVersion.ValueString()
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.RelationalDbms",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		nil,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
