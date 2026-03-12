package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &SearchEntityFunction{}

type SearchEntityFunction struct{}

func NewSearchEntityFunction() function.Function {
	return &SearchEntityFunction{}
}

func (f *SearchEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "search_entity"
}

func (f *SearchEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Search Entity blueprint component",
		Description: "Builds a Search Entity component with the correct type for use in a fractal's components list. If search_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Search Entity configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"search_id":    types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type searchEntityConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	SearchId    types.String `tfsdk:"search_id"`
}

func (f *SearchEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config searchEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	if !config.SearchId.IsNull() && !config.SearchId.IsUnknown() {
		deps = append(deps, config.SearchId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Storage.CaaS.SearchEntity",
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
