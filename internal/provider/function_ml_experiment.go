package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &MlExperimentFunction{}

type MlExperimentFunction struct{}

func NewMlExperimentFunction() function.Function {
	return &MlExperimentFunction{}
}

func (f *MlExperimentFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "ml_experiment"
}

func (f *MlExperimentFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates an ML Experiment blueprint component",
		Description: "Builds an ML Experiment component with the correct type for use in a fractal's components list. " +
			"If platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "ML Experiment configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                types.StringType,
					"display_name":      types.StringType,
					"description":       types.StringType,
					"version":           types.StringType,
					"platform_id":       types.StringType,
					"experiment_name":   types.StringType,
					"artifact_location": types.StringType,
				},
			},
		},
		Return: componentReturn(),
	}
}

type mlExperimentConfig struct {
	Id               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	Description      types.String `tfsdk:"description"`
	Version          types.String `tfsdk:"version"`
	PlatformId       types.String `tfsdk:"platform_id"`
	ExperimentName   types.String `tfsdk:"experiment_name"`
	ArtifactLocation types.String `tfsdk:"artifact_location"`
}

func (f *MlExperimentFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config mlExperimentConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ExperimentName.IsNull() && !config.ExperimentName.IsUnknown() {
		params["experimentName"] = config.ExperimentName.ValueString()
	}
	if !config.ArtifactLocation.IsNull() && !config.ArtifactLocation.IsUnknown() {
		params["artifactLocation"] = config.ArtifactLocation.ValueString()
	}

	var deps []string
	if !config.PlatformId.IsNull() && !config.PlatformId.IsUnknown() {
		deps = append(deps, config.PlatformId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.MlExperiment",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		params,
		deps,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
