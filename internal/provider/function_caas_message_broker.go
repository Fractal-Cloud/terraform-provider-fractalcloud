package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &CaaSMessageBrokerFunction{}

type CaaSMessageBrokerFunction struct{}

func NewCaaSMessageBrokerFunction() function.Function {
	return &CaaSMessageBrokerFunction{}
}

func (f *CaaSMessageBrokerFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "caas_message_broker"
}

func (f *CaaSMessageBrokerFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Containerized Message Broker blueprint component",
		Description: "Builds a Containerized Message Broker component with the correct type for use in a fractal's components list. If container_platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Containerized Message Broker configuration",
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

type caasMessageBrokerConfig struct {
	Id                  types.String `tfsdk:"id"`
	DisplayName         types.String `tfsdk:"display_name"`
	Description         types.String `tfsdk:"description"`
	Version             types.String `tfsdk:"version"`
	ContainerPlatformId types.String `tfsdk:"container_platform_id"`
}

func (f *CaaSMessageBrokerFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config caasMessageBrokerConfig
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
		"Messaging.CaaS.Broker",
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
