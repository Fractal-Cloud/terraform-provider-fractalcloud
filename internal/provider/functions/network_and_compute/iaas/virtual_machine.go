package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &VirtualMachineFunction{}

type VirtualMachineFunction struct{}

func NewVirtualMachineFunction() function.Function {
	return &VirtualMachineFunction{}
}

func (f *VirtualMachineFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "network_and_compute_iaas_virtual_machine"
}

func (f *VirtualMachineFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a VirtualMachine blueprint component",
		Description: "Builds a VirtualMachine (VM/EC2) component with the correct type for use in a fractal's components list. " +
			"If subnet is provided, it is validated as a Subnet and automatically added as a dependency. " +
			"Use links to define runtime relationships to other components, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "VirtualMachine configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"subnet":       components.ComponentObjectType,
					"links": types.ListType{
						ElemType: types.ObjectType{AttrTypes: components.GenericLinkAttrTypes},
					},
					"security_groups": types.ListType{ElemType: components.ComponentObjectType},
				},
			},
		},
		Return: components.ComponentReturn(),
	}
}

type virtualMachineConfig struct {
	Id             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Subnet         types.Object `tfsdk:"subnet"`
	Links          types.List   `tfsdk:"links"`
	SecurityGroups types.List   `tfsdk:"security_groups"`
}

func (f *VirtualMachineFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config virtualMachineConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	subnetId, funcErr := components.ExtractDependency(config.Subnet, "NetworkAndCompute.IaaS.Subnet")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if subnetId != "" {
		deps = append(deps, subnetId)
	}

	// Build links from generic links and SG memberships
	var links []components.ComponentLink

	if !config.Links.IsNull() && !config.Links.IsUnknown() {
		var genericLinks []components.GenericLinkConfig
		diags := config.Links.ElementsAs(ctx, &genericLinks, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse links")
			return
		}
		resolved, funcErr := components.GenericLinksToComponentLinks(genericLinks)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
			return
		}
		links = append(links, resolved...)
	}

	if !config.SecurityGroups.IsNull() && !config.SecurityGroups.IsUnknown() {
		var sgObjects []types.Object
		diags := config.SecurityGroups.ElementsAs(ctx, &sgObjects, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse security_groups")
			return
		}
		sgLinks, funcErr := components.SgMembershipLinks(sgObjects)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
			return
		}
		links = append(links, sgLinks...)
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.VirtualMachine",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		nil,
		deps,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
