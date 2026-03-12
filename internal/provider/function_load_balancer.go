package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &LoadBalancerFunction{}

type LoadBalancerFunction struct{}

func NewLoadBalancerFunction() function.Function {
	return &LoadBalancerFunction{}
}

func (f *LoadBalancerFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "load_balancer"
}

func (f *LoadBalancerFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a LoadBalancer blueprint component",
		Description: "Builds a LoadBalancer component with the correct type for use in a fractal's components list. " +
			"Use links to connect to backend workloads and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "LoadBalancer configuration",
				AttributeTypes: map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"version":      types.StringType,
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

type loadBalancerConfig struct {
	Id             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Version        types.String `tfsdk:"version"`
	Links          types.List   `tfsdk:"links"`
	SecurityGroups types.List   `tfsdk:"security_groups"`
}

func (f *LoadBalancerFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config loadBalancerConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
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
		"NetworkAndCompute.IaaS.LoadBalancer",
		optionalString(config.DisplayName),
		optionalString(config.Description),
		optionalString(config.Version),
		nil,
		nil,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
