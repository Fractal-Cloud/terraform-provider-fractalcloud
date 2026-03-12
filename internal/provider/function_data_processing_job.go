package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &DataProcessingJobFunction{}

type DataProcessingJobFunction struct{}

func NewDataProcessingJobFunction() function.Function {
	return &DataProcessingJobFunction{}
}

func (f *DataProcessingJobFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "data_processing_job"
}

func (f *DataProcessingJobFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Data Processing Job blueprint component",
		Description: "Builds a Data Processing Job component with the correct type for use in a fractal's components list. " +
			"If platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Data Processing Job configuration",
				AttributeTypes: map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"version":          types.StringType,
					"platform_id":      types.StringType,
					"job_name":         types.StringType,
					"task_type":        types.StringType,
					"notebook_path":    types.StringType,
					"python_file":      types.StringType,
					"main_class_name":  types.StringType,
					"jar_uri":          types.StringType,
					"cron_schedule":    types.StringType,
					"max_retries":      types.Int64Type,
					"existing_cluster": types.BoolType,
					"parameters":       types.ListType{ElemType: types.StringType},
				},
			},
		},
		Return: componentReturn(),
	}
}

type dataProcessingJobConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Version         types.String `tfsdk:"version"`
	PlatformId      types.String `tfsdk:"platform_id"`
	JobName         types.String `tfsdk:"job_name"`
	TaskType        types.String `tfsdk:"task_type"`
	NotebookPath    types.String `tfsdk:"notebook_path"`
	PythonFile      types.String `tfsdk:"python_file"`
	MainClassName   types.String `tfsdk:"main_class_name"`
	JarUri          types.String `tfsdk:"jar_uri"`
	CronSchedule    types.String `tfsdk:"cron_schedule"`
	MaxRetries      types.Int64  `tfsdk:"max_retries"`
	ExistingCluster types.Bool   `tfsdk:"existing_cluster"`
	Parameters      types.List   `tfsdk:"parameters"`
}

func (f *DataProcessingJobFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config dataProcessingJobConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.JobName.IsNull() && !config.JobName.IsUnknown() {
		params["jobName"] = config.JobName.ValueString()
	}
	if !config.TaskType.IsNull() && !config.TaskType.IsUnknown() {
		params["taskType"] = config.TaskType.ValueString()
	}
	if !config.NotebookPath.IsNull() && !config.NotebookPath.IsUnknown() {
		params["notebookPath"] = config.NotebookPath.ValueString()
	}
	if !config.PythonFile.IsNull() && !config.PythonFile.IsUnknown() {
		params["pythonFile"] = config.PythonFile.ValueString()
	}
	if !config.MainClassName.IsNull() && !config.MainClassName.IsUnknown() {
		params["mainClassName"] = config.MainClassName.ValueString()
	}
	if !config.JarUri.IsNull() && !config.JarUri.IsUnknown() {
		params["jarUri"] = config.JarUri.ValueString()
	}
	if !config.CronSchedule.IsNull() && !config.CronSchedule.IsUnknown() {
		params["cronSchedule"] = config.CronSchedule.ValueString()
	}
	if !config.MaxRetries.IsNull() && !config.MaxRetries.IsUnknown() {
		params["maxRetries"] = fmt.Sprintf("%d", config.MaxRetries.ValueInt64())
	}
	if !config.ExistingCluster.IsNull() && !config.ExistingCluster.IsUnknown() {
		params["existingCluster"] = fmt.Sprintf("%t", config.ExistingCluster.ValueBool())
	}

	if !config.Parameters.IsNull() && !config.Parameters.IsUnknown() {
		var jobParams []string
		diags := config.Parameters.ElementsAs(ctx, &jobParams, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse parameters")
			return
		}
		if len(jobParams) > 0 {
			b, err := json.Marshal(jobParams)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize parameters: %s", err))
				return
			}
			params["parameters"] = string(b)
		}
	}

	var deps []string
	if !config.PlatformId.IsNull() && !config.PlatformId.IsUnknown() {
		deps = append(deps, config.PlatformId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.DataProcessingJob",
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
