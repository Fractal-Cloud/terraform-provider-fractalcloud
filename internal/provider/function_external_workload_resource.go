package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &ExternalWorkloadResourceFunction{}

type ExternalWorkloadResourceFunction struct{}

func NewExternalWorkloadResourceFunction() function.Function {
	return &ExternalWorkloadResourceFunction{}
}

func (f *ExternalWorkloadResourceFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "external_workload_resource"
}

func (f *ExternalWorkloadResourceFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates an External Workload Resource blueprint component",
		Description: "Builds an External Workload Resource (unmanaged) component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "External Workload Resource configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type externalWorkloadResourceConfig struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
}

func (f *ExternalWorkloadResourceFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config externalWorkloadResourceConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"CustomWorkloads.SaaS.Unmanaged",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		nil,
		nil,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
