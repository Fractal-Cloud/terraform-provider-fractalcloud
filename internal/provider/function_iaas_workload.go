package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &IaaSWorkloadFunction{}

type IaaSWorkloadFunction struct{}

func NewIaaSWorkloadFunction() function.Function {
	return &IaaSWorkloadFunction{}
}

func (f *IaaSWorkloadFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "iaas_workload"
}

func (f *IaaSWorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates an IaaS Workload blueprint component",
		Description: "Builds an IaaS Workload component (running on VMs or bare metal) with the correct type for use in a fractal's components list. " +
			"If vm_id or subnet_id are provided, they are automatically added as dependencies. " +
			"Use links for port-based traffic rules to other workloads, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "IaaS Workload configuration",
				AttributeTypes: map[string]attr.Type{
					"id":              types.StringType,
					"display_name":    types.StringType,
					"description":     types.StringType,
					"version":         types.StringType,
					"container_image": types.StringType,
					"container_port":  types.Int64Type,
					"container_name":  types.StringType,
					"cpu":             types.StringType,
					"memory":          types.StringType,
					"desired_count":   types.Int64Type,
					"vm_id":           types.StringType,
					"subnet_id":       types.StringType,
					"links": types.ListType{
						ElemType: types.ObjectType{AttrTypes: portLinkAttrTypes},
					},
					"security_groups": types.ListType{ElemType: types.StringType},
				},
			},
		},
		Return: componentReturn(),
	}
}

type iaasWorkloadConfig struct {
	Id             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Version        types.String `tfsdk:"version"`
	ContainerImage types.String `tfsdk:"container_image"`
	ContainerPort  types.Int64  `tfsdk:"container_port"`
	ContainerName  types.String `tfsdk:"container_name"`
	Cpu            types.String `tfsdk:"cpu"`
	Memory         types.String `tfsdk:"memory"`
	DesiredCount   types.Int64  `tfsdk:"desired_count"`
	VmId           types.String `tfsdk:"vm_id"`
	SubnetId       types.String `tfsdk:"subnet_id"`
	Links          types.List   `tfsdk:"links"`
	SecurityGroups types.List   `tfsdk:"security_groups"`
}

func (f *IaaSWorkloadFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config iaasWorkloadConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ContainerImage.IsNull() && !config.ContainerImage.IsUnknown() {
		params["containerImage"] = config.ContainerImage.ValueString()
	}
	if !config.ContainerPort.IsNull() && !config.ContainerPort.IsUnknown() {
		params["containerPort"] = fmt.Sprintf("%d", config.ContainerPort.ValueInt64())
	}
	if !config.ContainerName.IsNull() && !config.ContainerName.IsUnknown() {
		params["containerName"] = config.ContainerName.ValueString()
	}
	if !config.Cpu.IsNull() && !config.Cpu.IsUnknown() {
		params["cpu"] = config.Cpu.ValueString()
	}
	if !config.Memory.IsNull() && !config.Memory.IsUnknown() {
		params["memory"] = config.Memory.ValueString()
	}
	if !config.DesiredCount.IsNull() && !config.DesiredCount.IsUnknown() {
		params["desiredCount"] = fmt.Sprintf("%d", config.DesiredCount.ValueInt64())
	}

	var deps []string
	if !config.VmId.IsNull() && !config.VmId.IsUnknown() {
		deps = append(deps, config.VmId.ValueString())
	}
	if !config.SubnetId.IsNull() && !config.SubnetId.IsUnknown() {
		deps = append(deps, config.SubnetId.ValueString())
	}

	var links []componentLink

	if !config.Links.IsNull() && !config.Links.IsUnknown() {
		var portLinks []portLinkConfig
		diags := config.Links.ElementsAs(ctx, &portLinks, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse links")
			return
		}
		links = append(links, portLinksToComponentLinks(portLinks)...)
	}

	if !config.SecurityGroups.IsNull() && !config.SecurityGroups.IsUnknown() {
		var sgIds []string
		diags := config.SecurityGroups.ElementsAs(ctx, &sgIds, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse security_groups")
			return
		}
		links = append(links, sgMembershipLinks(sgIds)...)
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"CustomWorkloads.IaaS.Workload",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		deps,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
