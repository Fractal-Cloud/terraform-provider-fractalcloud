package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &ComputeClusterFunction{}

type ComputeClusterFunction struct{}

func NewComputeClusterFunction() function.Function {
	return &ComputeClusterFunction{}
}

func (f *ComputeClusterFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "compute_cluster"
}

func (f *ComputeClusterFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Compute Cluster blueprint component",
		Description: "Builds a Compute Cluster (managed Spark cluster) component with the correct type for use in a fractal's components list. " +
			"If platform_id is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Compute Cluster configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                       types.StringType,
					"display_name":             types.StringType,
					"description":              types.StringType,
					"version":                  types.StringType,
					"platform_id":              types.StringType,
					"cluster_name":             types.StringType,
					"spark_version":            types.StringType,
					"node_type_id":             types.StringType,
					"num_workers":              types.Int64Type,
					"min_workers":              types.Int64Type,
					"max_workers":              types.Int64Type,
					"auto_termination_minutes": types.Int64Type,
					"spark_conf":               types.MapType{ElemType: types.StringType},
					"pypi_libraries":           types.ListType{ElemType: types.StringType},
					"maven_libraries":          types.ListType{ElemType: types.StringType},
				},
			},
		},
		Return: componentReturn(),
	}
}

type computeClusterConfig struct {
	Id                     types.String `tfsdk:"id"`
	DisplayName            types.String `tfsdk:"display_name"`
	Description            types.String `tfsdk:"description"`
	Version                types.String `tfsdk:"version"`
	PlatformId             types.String `tfsdk:"platform_id"`
	ClusterName            types.String `tfsdk:"cluster_name"`
	SparkVersion           types.String `tfsdk:"spark_version"`
	NodeTypeId             types.String `tfsdk:"node_type_id"`
	NumWorkers             types.Int64  `tfsdk:"num_workers"`
	MinWorkers             types.Int64  `tfsdk:"min_workers"`
	MaxWorkers             types.Int64  `tfsdk:"max_workers"`
	AutoTerminationMinutes types.Int64  `tfsdk:"auto_termination_minutes"`
	SparkConf              types.Map    `tfsdk:"spark_conf"`
	PypiLibraries          types.List   `tfsdk:"pypi_libraries"`
	MavenLibraries         types.List   `tfsdk:"maven_libraries"`
}

func (f *ComputeClusterFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config computeClusterConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ClusterName.IsNull() && !config.ClusterName.IsUnknown() {
		params["clusterName"] = config.ClusterName.ValueString()
	}
	if !config.SparkVersion.IsNull() && !config.SparkVersion.IsUnknown() {
		params["sparkVersion"] = config.SparkVersion.ValueString()
	}
	if !config.NodeTypeId.IsNull() && !config.NodeTypeId.IsUnknown() {
		params["nodeTypeId"] = config.NodeTypeId.ValueString()
	}
	if !config.NumWorkers.IsNull() && !config.NumWorkers.IsUnknown() {
		params["numWorkers"] = fmt.Sprintf("%d", config.NumWorkers.ValueInt64())
	}
	if !config.MinWorkers.IsNull() && !config.MinWorkers.IsUnknown() {
		params["minWorkers"] = fmt.Sprintf("%d", config.MinWorkers.ValueInt64())
	}
	if !config.MaxWorkers.IsNull() && !config.MaxWorkers.IsUnknown() {
		params["maxWorkers"] = fmt.Sprintf("%d", config.MaxWorkers.ValueInt64())
	}
	if !config.AutoTerminationMinutes.IsNull() && !config.AutoTerminationMinutes.IsUnknown() {
		params["autoTerminationMinutes"] = fmt.Sprintf("%d", config.AutoTerminationMinutes.ValueInt64())
	}

	if !config.SparkConf.IsNull() && !config.SparkConf.IsUnknown() {
		confMap := make(map[string]string)
		diags := config.SparkConf.ElementsAs(ctx, &confMap, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse spark_conf")
			return
		}
		if len(confMap) > 0 {
			b, err := json.Marshal(confMap)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize spark_conf: %s", err))
				return
			}
			params["sparkConf"] = string(b)
		}
	}

	if !config.PypiLibraries.IsNull() && !config.PypiLibraries.IsUnknown() {
		var libs []string
		diags := config.PypiLibraries.ElementsAs(ctx, &libs, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse pypi_libraries")
			return
		}
		if len(libs) > 0 {
			b, err := json.Marshal(libs)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize pypi_libraries: %s", err))
				return
			}
			params["pypiLibraries"] = string(b)
		}
	}

	if !config.MavenLibraries.IsNull() && !config.MavenLibraries.IsUnknown() {
		var libs []string
		diags := config.MavenLibraries.ElementsAs(ctx, &libs, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse maven_libraries")
			return
		}
		if len(libs) > 0 {
			b, err := json.Marshal(libs)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize maven_libraries: %s", err))
				return
			}
			params["mavenLibraries"] = string(b)
		}
	}

	var deps []string
	if !config.PlatformId.IsNull() && !config.PlatformId.IsUnknown() {
		deps = append(deps, config.PlatformId.ValueString())
	}

	result, funcErr := buildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.ComputeCluster",
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
