package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &CaaSMessageBrokerEntityFunction{}

type CaaSMessageBrokerEntityFunction struct{}

func NewCaaSMessageBrokerEntityFunction() function.Function {
	return &CaaSMessageBrokerEntityFunction{}
}

func (f *CaaSMessageBrokerEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "caas_message_broker_entity"
}

func (f *CaaSMessageBrokerEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Containerized Message Broker Entity blueprint component",
		Description: "Builds a Containerized Message Broker Entity component with the correct type for use in a fractal's components list. If broker_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Containerized Message Broker Entity configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"broker_id":    types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type caasMessageBrokerEntityConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	BrokerId    types.String `tfsdk:"broker_id"`
}

func (f *CaaSMessageBrokerEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config caasMessageBrokerEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	if !config.BrokerId.IsNull() && !config.BrokerId.IsUnknown() {
		deps = append(deps, config.BrokerId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Messaging.CaaS.Entity",
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
