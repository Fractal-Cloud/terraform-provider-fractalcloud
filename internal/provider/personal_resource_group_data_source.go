package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &PersonalResourceGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &PersonalResourceGroupsDataSource{}
)

// NewPersonalResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewPersonalResourceGroupDataSource() datasource.DataSource {
	return &PersonalResourceGroupsDataSource{}
}

// PersonalResourceGroupsDataSource is the data source implementation.
type PersonalResourceGroupsDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *PersonalResourceGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PersonalResourceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_resource_group"
}

// Schema defines the schema for the data source.
func (d *PersonalResourceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"type":       basetypes.StringType{},
					"owner_id":   basetypes.StringType{},
					"short_name": basetypes.StringType{},
				},
			},
			"short_name": schema.StringAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"icon": schema.StringAttribute{
				Computed: true,
			},
			"live_systems_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"fractals_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// PersonalResourceGroupModel maps resource group schema data.
type PersonalResourceGroupModel struct {
	Id             types.Object `tfsdk:"id"`
	ShortName      types.String `tfsdk:"short_name"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Status         types.String `tfsdk:"status"`
	Icon           types.String `tfsdk:"icon"`
	LiveSystemsIds types.List   `tfsdk:"live_systems_ids"`
	FractalsIds    types.List   `tfsdk:"fractals_ids"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

// Read refreshes the Terraform state with the latest data.
func (d *PersonalResourceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PersonalResourceGroupModel

	// Read user configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required input isn't unknown/null
	if config.ShortName.IsUnknown() || config.ShortName.IsNull() {
		resp.Diagnostics.AddError(
			"Unknown required value",
			fmt.Sprintf("ShortName is required, but its value is unknown, null or empty."),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	var resourceGroupId = fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: config.ShortName.ValueString(),
	}

	resourceGroup, err := d.client.GetPersonalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not read Personal Fractal Cloud Resource Group with ShortName "+config.ShortName.ValueString()+": "+err.Error())
		return
	}

	if resourceGroup == nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not find Fractal Cloud Resource Group Id "+config.ShortName.ValueString())
		return
	}

	fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
	resp.Diagnostics.Append(diags...)

	liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
	resp.Diagnostics.Append(diags...)

	idTypes := map[string]attr.Type{
		"type":       types.StringType,
		"owner_id":   types.StringType,
		"short_name": types.StringType,
	}

	// Build state
	state := PersonalResourceGroupModel{
		Id: types.ObjectValueMust(idTypes, map[string]attr.Value{
			"type":       types.StringValue(resourceGroup.Id.Type),
			"owner_id":   types.StringValue(resourceGroup.Id.OwnerId),
			"short_name": types.StringValue(resourceGroup.Id.ShortName),
		}),
		ShortName:      types.StringValue(resourceGroup.Id.ShortName),
		DisplayName:    types.StringValue(resourceGroup.DisplayName),
		Description:    types.StringValue(resourceGroup.Description),
		Icon:           types.StringValue(resourceGroup.Icon),
		FractalsIds:    fractalsIds,
		LiveSystemsIds: liveSystemsIds,
		CreatedAt:      types.StringValue(resourceGroup.CreatedAt),
		UpdatedAt:      types.StringValue(resourceGroup.UpdatedAt),
	}

	// Write state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
