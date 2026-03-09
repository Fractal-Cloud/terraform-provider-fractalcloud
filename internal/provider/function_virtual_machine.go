package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &VirtualMachineFunction{}

type VirtualMachineFunction struct{}

func NewVirtualMachineFunction() function.Function {
	return &VirtualMachineFunction{}
}

func (f *VirtualMachineFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "virtual_machine"
}

func (f *VirtualMachineFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a VirtualMachine blueprint component",
		Description: "Builds a VirtualMachine (VM/EC2) component with the correct type for use in a fractal's components list. If subnet_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "VirtualMachine configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
					"subnet_id":    types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type virtualMachineConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	SubnetId    types.String `tfsdk:"subnet_id"`
}

func (f *VirtualMachineFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config virtualMachineConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	if !config.SubnetId.IsNull() && !config.SubnetId.IsUnknown() {
		deps = append(deps, config.SubnetId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.VirtualMachine",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		nil,
		deps,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
