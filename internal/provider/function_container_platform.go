package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &ContainerPlatformFunction{}

type ContainerPlatformFunction struct{}

func NewContainerPlatformFunction() function.Function {
	return &ContainerPlatformFunction{}
}

func (f *ContainerPlatformFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "container_platform"
}

func (f *ContainerPlatformFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a ContainerPlatform blueprint component",
		Description: "Builds a ContainerPlatform (managed Kubernetes) component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "ContainerPlatform configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type containerPlatformConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
}

func (f *ContainerPlatformFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config containerPlatformConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.PaaS.ContainerPlatform",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		nil,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
