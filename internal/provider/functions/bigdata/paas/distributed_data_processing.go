package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &BigdataPaasDistributedDataProcessingFunction{}

type BigdataPaasDistributedDataProcessingFunction struct{}

func NewBigdataPaasDistributedDataProcessingFunction() function.Function {
	return &BigdataPaasDistributedDataProcessingFunction{}
}

func (f *BigdataPaasDistributedDataProcessingFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_paas_distributed_data_processing"
}

func (f *BigdataPaasDistributedDataProcessingFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a BigData PaaS Distributed Data Processing blueprint component",
		Description: "Builds a BigData PaaS Distributed Data Processing component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Distributed Data Processing configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"pricing_tier": types.StringType,
				},
			},
		},
		Return: components.ComponentReturn(),
	}
}

type bigdataPaasDistributedDataProcessingConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	PricingTier types.String `tfsdk:"pricing_tier"`
}

func (f *BigdataPaasDistributedDataProcessingFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataPaasDistributedDataProcessingConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}
	if !config.PricingTier.IsNull() && !config.PricingTier.IsUnknown() {
		params["pricingTier"] = config.PricingTier.ValueString()
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.DistributedDataProcessing",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
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
