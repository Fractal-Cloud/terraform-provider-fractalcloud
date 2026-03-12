package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &DistributedDataProcessingFunction{}

type DistributedDataProcessingFunction struct{}

func NewDistributedDataProcessingFunction() function.Function {
	return &DistributedDataProcessingFunction{}
}

func (f *DistributedDataProcessingFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "distributed_data_processing"
}

func (f *DistributedDataProcessingFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Distributed Data Processing blueprint component",
		Description: "Builds a Distributed Data Processing platform component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Distributed Data Processing configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"pricing_tier": types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type distributedDataProcessingConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	PricingTier types.String `tfsdk:"pricing_tier"`
}

func (f *DistributedDataProcessingFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config distributedDataProcessingConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.PricingTier.IsNull() && !config.PricingTier.IsUnknown() {
		params["pricingTier"] = config.PricingTier.ValueString()
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.DistributedDataProcessing",
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
