package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &LoggingFunction{}

type LoggingFunction struct{}

func NewLoggingFunction() function.Function {
	return &LoggingFunction{}
}

func (f *LoggingFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "logging"
}

func (f *LoggingFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates an Observability Logging blueprint component",
		Description: "Builds an Observability Logging component with the correct type for use in a fractal's components list. If container_platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Logging configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                    types.StringType,
					"display_name":          types.StringType,
					"description":           types.StringType,
					"version":               types.StringType,
					"container_platform_id": types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type loggingConfig struct {
	Id                  types.String `tfsdk:"id"`
	DisplayName         types.String `tfsdk:"display_name"`
	Description         types.String `tfsdk:"description"`
	Version             types.String `tfsdk:"version"`
	ContainerPlatformId types.String `tfsdk:"container_platform_id"`
}

func (f *LoggingFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config loggingConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	if !config.ContainerPlatformId.IsNull() && !config.ContainerPlatformId.IsUnknown() {
		deps = append(deps, config.ContainerPlatformId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Observability.CaaS.Logging",
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
