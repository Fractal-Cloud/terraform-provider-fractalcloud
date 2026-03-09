package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &VirtualNetworkFunction{}

type VirtualNetworkFunction struct{}

func NewVirtualNetworkFunction() function.Function {
	return &VirtualNetworkFunction{}
}

func (f *VirtualNetworkFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "virtual_network"
}

func (f *VirtualNetworkFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a VirtualNetwork blueprint component",
		Description: "Builds a VirtualNetwork (VPC) component with the correct type and parameters for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "VirtualNetwork configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"cidr_block":   types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type virtualNetworkConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	CidrBlock   types.String `tfsdk:"cidr_block"`
}

func (f *VirtualNetworkFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config virtualNetworkConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}
	if !config.CidrBlock.IsNull() && !config.CidrBlock.IsUnknown() {
		params["cidrBlock"] = config.CidrBlock.ValueString()
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.VirtualNetwork",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
