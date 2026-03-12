package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &MessageBrokerEntityFunction{}

type MessageBrokerEntityFunction struct{}

func NewMessageBrokerEntityFunction() function.Function {
	return &MessageBrokerEntityFunction{}
}

func (f *MessageBrokerEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "message_broker_entity"
}

func (f *MessageBrokerEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Messaging Platform Entity blueprint component",
		Description: "Builds a Messaging Platform Entity component with the correct type for use in a fractal's components list. " +
			"If broker_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Messaging Platform Entity configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                       types.StringType,
					"display_name":             types.StringType,
					"description":              types.StringType,
					"version":                  types.StringType,
					"broker_id":                types.StringType,
					"message_retention_hours":  types.Int64Type,
				},
			},
		},
		Return: componentReturn(),
	}
}

type messageBrokerEntityConfig struct {
	Id                    types.String `tfsdk:"id"`
	DisplayName           types.String `tfsdk:"display_name"`
	Description           types.String `tfsdk:"description"`
	Version               types.String `tfsdk:"version"`
	BrokerId              types.String `tfsdk:"broker_id"`
	MessageRetentionHours types.Int64  `tfsdk:"message_retention_hours"`
}

func (f *MessageBrokerEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config messageBrokerEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.MessageRetentionHours.IsNull() && !config.MessageRetentionHours.IsUnknown() {
		params["messageRetentionHours"] = fmt.Sprintf("%d", config.MessageRetentionHours.ValueInt64())
	}

	var deps []string
	if !config.BrokerId.IsNull() && !config.BrokerId.IsUnknown() {
		deps = append(deps, config.BrokerId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"Messaging.PaaS.Entity",
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
