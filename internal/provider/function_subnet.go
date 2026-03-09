package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &SubnetFunction{}

type SubnetFunction struct{}

func NewSubnetFunction() function.Function {
	return &SubnetFunction{}
}

func (f *SubnetFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "subnet"
}

func (f *SubnetFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Subnet blueprint component",
		Description: "Builds a Subnet component with the correct type and parameters for use in a fractal's components list. If vpc_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Subnet configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                types.StringType,
					"display_name":      types.StringType,
					"description":       types.StringType,
					"version":           types.StringType,
					"cidr_block":        types.StringType,
					"availability_zone": types.StringType,
					"vpc_id":            types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type subnetConfig struct {
	Id               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	Description      types.String `tfsdk:"description"`
	Version          types.String `tfsdk:"version"`
	CidrBlock        types.String `tfsdk:"cidr_block"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	VpcId            types.String `tfsdk:"vpc_id"`
}

func (f *SubnetFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config subnetConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}
	if !config.CidrBlock.IsNull() && !config.CidrBlock.IsUnknown() {
		params["cidrBlock"] = config.CidrBlock.ValueString()
	}
	if !config.AvailabilityZone.IsNull() && !config.AvailabilityZone.IsUnknown() {
		params["availabilityZone"] = config.AvailabilityZone.ValueString()
	}

	var deps []string
	if !config.VpcId.IsNull() && !config.VpcId.IsUnknown() {
		deps = append(deps, config.VpcId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.Subnet",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		deps,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
