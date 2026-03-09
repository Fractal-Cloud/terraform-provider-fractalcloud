package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &WorkloadFunction{}

type WorkloadFunction struct{}

func NewWorkloadFunction() function.Function {
	return &WorkloadFunction{}
}

func (f *WorkloadFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "workload"
}

func (f *WorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Workload blueprint component",
		Description: "Builds a containerized Workload component with the correct type and parameters for use in a fractal's components list. If platform_id or subnet_id are provided, they are automatically added as dependencies.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Workload configuration",
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
					"platform_id":     types.StringType,
					"subnet_id":       types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type workloadConfig struct {
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
	PlatformId     types.String `tfsdk:"platform_id"`
	SubnetId       types.String `tfsdk:"subnet_id"`
}

func (f *WorkloadFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config workloadConfig
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
	if !config.PlatformId.IsNull() && !config.PlatformId.IsUnknown() {
		deps = append(deps, config.PlatformId.ValueString())
	}
	if !config.SubnetId.IsNull() && !config.SubnetId.IsUnknown() {
		deps = append(deps, config.SubnetId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"CustomWorkloads.CaaS.Workload",
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
