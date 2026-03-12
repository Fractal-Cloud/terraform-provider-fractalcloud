package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &ServiceMeshSecurityFunction{}

type ServiceMeshSecurityFunction struct{}

func NewServiceMeshSecurityFunction() function.Function {
	return &ServiceMeshSecurityFunction{}
}

func (f *ServiceMeshSecurityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "service_mesh_security"
}

func (f *ServiceMeshSecurityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Service Mesh Security blueprint component",
		Description: "Builds a Service Mesh Security component with the correct type for use in a fractal's components list. If container_platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Service Mesh Security configuration",
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

type serviceMeshSecurityConfig struct {
	Id                  types.String `tfsdk:"id"`
	DisplayName         types.String `tfsdk:"display_name"`
	Description         types.String `tfsdk:"description"`
	Version             types.String `tfsdk:"version"`
	ContainerPlatformId types.String `tfsdk:"container_platform_id"`
}

func (f *ServiceMeshSecurityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config serviceMeshSecurityConfig
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
		"Security.CaaS.ServiceMeshsecurity",
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
