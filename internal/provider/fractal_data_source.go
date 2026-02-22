package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &FractalsDataSource{}
	_ datasource.DataSourceWithConfigure = &FractalsDataSource{}
)

// NewFractalDataSource is a helper function to simplify the provider implementation.
func NewFractalDataSource() datasource.DataSource {
	return &FractalsDataSource{}
}

// FractalsDataSource is the data source implementation.
type FractalsDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *FractalsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fractalCloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fractalCloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *FractalsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fractal"
}

// Schema defines the schema for the data source.
func (d *FractalsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_group_id": schema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"type":       basetypes.StringType{},
					"owner_id":   basetypes.StringType{},
					"short_name": basetypes.StringType{},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"version": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"is_private": schema.BoolAttribute{
				Computed: true,
			},
			"components": schema.ListAttribute{
				Computed: true,
				ElementType: basetypes.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":                  basetypes.StringType{},
						"type":                basetypes.StringType{},
						"display_name":        basetypes.StringType{},
						"description":         basetypes.StringType{},
						"version":             basetypes.StringType{},
						"is_locked":           basetypes.BoolType{},
						"recreate_on_failure": basetypes.BoolType{},
						"parameters": basetypes.MapType{
							ElemType: basetypes.StringType{},
						},
						"dependencies_ids": basetypes.ListType{
							ElemType: basetypes.StringType{},
						},
						"links": basetypes.ListType{
							ElemType: basetypes.ObjectType{
								AttrTypes: map[string]attr.Type{
									"component_id": basetypes.StringType{},
									"settings": basetypes.MapType{
										ElemType: basetypes.StringType{},
									},
								},
							},
						},
						"output_fields": basetypes.ListType{
							ElemType: basetypes.StringType{},
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type LinkModel struct {
	ComponentId types.String `tfsdk:"component_id"`
	Settings    types.Map    `tfsdk:"settings"`
}

type ComponentModel struct {
	Id                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	Version           types.String `tfsdk:"version"`
	IsLocked          types.Bool   `tfsdk:"is_locked"`
	RecreateOnFailure types.Bool   `tfsdk:"recreate_on_failure"`
	Parameters        types.Map    `tfsdk:"parameters"`
	DependenciesIds   types.List   `tfsdk:"dependencies_ids"`
	Links             types.List   `tfsdk:"links"`
	OutputFields      types.List   `tfsdk:"output_fields"`
}

type ResourceGroupIdModel struct {
	Type      types.String `tfsdk:"type"`
	OwnerId   types.String `tfsdk:"owner_id"`
	ShortName types.String `tfsdk:"short_name"`
}

// BlueprintModel maps resource group schema data.
type BlueprintModel struct {
	ResourceGroupId ResourceGroupIdModel `tfsdk:"resource_group_id"`
	Name            types.String         `tfsdk:"name"`
	Version         types.String         `tfsdk:"version"`
	Description     types.String         `tfsdk:"description"`
	IsPrivate       types.Bool           `tfsdk:"is_private"`
	Components      types.List           `tfsdk:"components"`
	CreatedAt       types.String         `tfsdk:"created_at"`
}

// Read refreshes the Terraform state with the latest data.
func (d *FractalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config BlueprintModel

	// Read user configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, done := GetFractalModel(ctx, config, d.client, resp.Diagnostics)
	if done {
		return
	}

	// Write state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func GetFractalModel(
	ctx context.Context,
	config BlueprintModel,
	client *fractalCloud.Client,
	diagnostics diag.Diagnostics) (BlueprintModel, bool) {
	var fractalId = fractalCloud.FractalId{
		ResourceGroupId: fractalCloud.ResourceGroupId{
			Type:      config.ResourceGroupId.Type.ValueString(),
			OwnerId:   config.ResourceGroupId.OwnerId.ValueString(),
			ShortName: config.ResourceGroupId.ShortName.ValueString(),
		},
		Name:    config.Name.ValueString(),
		Version: config.Version.ValueString(),
	}

	blueprint, err := client.GetBlueprint(fractalId)
	if err != nil {
		diagnostics.AddError(
			"Error Reading Fractal Cloud Blueprint",
			"Could not read Fractal with Id "+fractalId.ToString()+": "+err.Error())
		return BlueprintModel{}, true
	}

	if blueprint == nil {
		diagnostics.AddError(
			"Error Reading Fractal Cloud Blueprint",
			"Could not find Fractal Cloud Blueprint with Id "+fractalId.ToString())
		return BlueprintModel{}, true
	}

	components := make([]ComponentModel, len(blueprint.Components))
	for i, component := range blueprint.Components {
		parameters, diags := types.MapValueFrom(ctx, types.StringType, component.Parameters)
		diagnostics.Append(diags...)

		dependenciesIds, diags := types.ListValueFrom(ctx, types.StringType, component.DependenciesIds)
		diagnostics.Append(diags...)

		links := make([]LinkModel, len(component.Links))
		for j, link := range component.Links {
			settings, diags := types.MapValueFrom(ctx, types.StringType, link.Settings)
			diagnostics.Append(diags...)

			links[j] = LinkModel{
				ComponentId: types.StringValue(link.ComponentId),
				Settings:    settings,
			}
		}

		linksToMap, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"component_id": types.StringType,
				"settings": types.MapType{
					ElemType: basetypes.StringType{},
				},
			},
		}, links)
		diagnostics.Append(diags...)

		outputFields, diags := types.ListValueFrom(ctx, types.StringType, component.OutputFields)
		diagnostics.Append(diags...)

		components[i] = ComponentModel{
			Id:                types.StringValue(component.Id),
			Type:              types.StringValue(component.Type),
			DisplayName:       types.StringValue(component.DisplayName),
			Description:       types.StringValue(component.Description),
			Version:           types.StringValue(component.Version),
			IsLocked:          types.BoolValue(component.IsLocked),
			RecreateOnFailure: types.BoolValue(component.RecreateOnFailure),
			Parameters:        parameters,
			DependenciesIds:   dependenciesIds,
			Links:             linksToMap,
			OutputFields:      outputFields,
		}
	}

	componentsToMap, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                  basetypes.StringType{},
			"type":                basetypes.StringType{},
			"display_name":        basetypes.StringType{},
			"description":         basetypes.StringType{},
			"version":             basetypes.StringType{},
			"is_locked":           basetypes.BoolType{},
			"recreate_on_failure": basetypes.BoolType{},
			"parameters": basetypes.MapType{
				ElemType: basetypes.StringType{},
			},
			"dependencies_ids": basetypes.ListType{
				ElemType: basetypes.StringType{},
			},
			"links": basetypes.ListType{
				ElemType: basetypes.ObjectType{
					AttrTypes: map[string]attr.Type{
						"component_id": basetypes.StringType{},
						"settings": basetypes.MapType{
							ElemType: basetypes.StringType{},
						},
					},
				},
			},
			"output_fields": basetypes.ListType{
				ElemType: basetypes.StringType{},
			},
		},
	}, components)
	diagnostics.Append(diags...)

	// Build state
	state := BlueprintModel{
		ResourceGroupId: config.ResourceGroupId,
		Name:            config.Name,
		Version:         config.Version,
		Description:     types.StringValue(blueprint.Description),
		IsPrivate:       types.BoolValue(blueprint.IsPrivate),
		Components:      componentsToMap,
		CreatedAt:       types.StringValue(blueprint.CreatedAt),
	}
	return state, false
}
