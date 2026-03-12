package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &FaaSWorkloadFunction{}

type FaaSWorkloadFunction struct{}

func NewFaaSWorkloadFunction() function.Function {
	return &FaaSWorkloadFunction{}
}

func (f *FaaSWorkloadFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "faas_workload"
}

func (f *FaaSWorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a FaaS Workload blueprint component",
		Description: "Builds a FaaS Workload (serverless function) component with the correct type for use in a fractal's components list. " +
			"If platform_id or subnet_id are provided, they are automatically added as dependencies. " +
			"Use links for port-based traffic rules to other workloads, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "FaaS Workload configuration",
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
					"runtime":         types.StringType,
					"memory_mb":       types.Int64Type,
					"timeout_seconds": types.Int64Type,
					"handler":         types.StringType,
					"platform_id":     types.StringType,
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

type faasWorkloadConfig struct {
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
	Runtime        types.String `tfsdk:"runtime"`
	MemoryMb       types.Int64  `tfsdk:"memory_mb"`
	TimeoutSeconds types.Int64  `tfsdk:"timeout_seconds"`
	Handler        types.String `tfsdk:"handler"`
	PlatformId     types.String `tfsdk:"platform_id"`
	SubnetId       types.String `tfsdk:"subnet_id"`
	Links          types.List   `tfsdk:"links"`
	SecurityGroups types.List   `tfsdk:"security_groups"`
}

func (f *FaaSWorkloadFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config faasWorkloadConfig
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
	if !config.Runtime.IsNull() && !config.Runtime.IsUnknown() {
		params["runtime"] = config.Runtime.ValueString()
	}
	if !config.MemoryMb.IsNull() && !config.MemoryMb.IsUnknown() {
		params["memoryMb"] = fmt.Sprintf("%d", config.MemoryMb.ValueInt64())
	}
	if !config.TimeoutSeconds.IsNull() && !config.TimeoutSeconds.IsUnknown() {
		params["timeoutSeconds"] = fmt.Sprintf("%d", config.TimeoutSeconds.ValueInt64())
	}
	if !config.Handler.IsNull() && !config.Handler.IsUnknown() {
		params["handler"] = config.Handler.ValueString()
	}

	var deps []string
	if !config.PlatformId.IsNull() && !config.PlatformId.IsUnknown() {
		deps = append(deps, config.PlatformId.ValueString())
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
		"CustomWorkloads.FaaS.Workload",
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
