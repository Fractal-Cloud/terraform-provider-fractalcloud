package caas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &StorageCaasSearchFunction{}

type StorageCaasSearchFunction struct{}

func NewStorageCaasSearchFunction() function.Function {
	return &StorageCaasSearchFunction{}
}

func (f *StorageCaasSearchFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "storage_caas_search"
}

func (f *StorageCaasSearchFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Containerized Search Platform blueprint component",
		Description: "Builds a Containerized Search Platform component with the correct type for use in a fractal's components list. " +
			"If container_platform is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Search Platform configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                 types.StringType,
					"display_name":       types.StringType,
					"description":        types.StringType,
					"version":            types.StringType,
					"container_platform": components.ComponentObjectType,
				},
			},
		},
		Return: components.ComponentReturn(),
	}
}

type storageCaasSearchConfig struct {
	Id                types.String `tfsdk:"id"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	Version           types.String `tfsdk:"version"`
	ContainerPlatform types.Object `tfsdk:"container_platform"`
}

func (f *StorageCaasSearchFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config storageCaasSearchConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	cpId, funcErr := components.ExtractDependency(config.ContainerPlatform, "NetworkAndCompute.PaaS.ContainerPlatform")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if cpId != "" {
		deps = append(deps, cpId)
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Storage.CaaS.Search",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		components.OptionalString(config.Version),
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
