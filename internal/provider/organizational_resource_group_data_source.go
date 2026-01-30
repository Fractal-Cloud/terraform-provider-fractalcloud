package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &OrganizationalResourceGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationalResourceGroupsDataSource{}
)

// NewOrganizationalResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewOrganizationalResourceGroupDataSource() datasource.DataSource {
	return &OrganizationalResourceGroupsDataSource{}
}

// OrganizationalResourceGroupsDataSource is the data source implementation.
type OrganizationalResourceGroupsDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationalResourceGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationalResourceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_resource_group"
}

// Schema defines the schema for the data source.
func (d *OrganizationalResourceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"short_name": schema.StringAttribute{
				Required: true,
			},
			"organization_id": schema.StringAttribute{
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
			"members_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"teams_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"managers_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
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
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_by": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// OrganizationalResourceGroupModel maps resource group schema data.
type OrganizationalResourceGroupModel struct {
	ShortName      types.String `tfsdk:"short_name"`
	OrganizationId types.String `tfsdk:"organization_id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Status         types.String `tfsdk:"status"`
	Icon           types.String `tfsdk:"icon"`
	MembersIds     types.List   `tfsdk:"members_ids"`
	TeamsIds       types.List   `tfsdk:"teams_ids"`
	ManagersIds    types.List   `tfsdk:"managers_ids"`
	LiveSystemsIds types.List   `tfsdk:"live_systems_ids"`
	FractalsIds    types.List   `tfsdk:"fractals_ids"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
	UpdatedBy      types.String `tfsdk:"updated_by"`
}

// Read refreshes the Terraform state with the latest data.
func (d *OrganizationalResourceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OrganizationalResourceGroupModel

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
		Type:      "Organizational",
		ShortName: config.ShortName.ValueString(),
	}

	resourceGroup, err := d.client.GetOrganizationalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not read Organizational Fractal Cloud Resource Group with ShortName "+config.ShortName.ValueString()+": "+err.Error())
		return
	}

	if resourceGroup == nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not find Fractal Cloud Resource Group ID "+config.ShortName.ValueString())
		return
	}

	membersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.MembersIds)
	resp.Diagnostics.Append(diags...)

	teamsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.TeamsIds)
	resp.Diagnostics.Append(diags...)

	managersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.ManagersIds)
	resp.Diagnostics.Append(diags...)

	fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
	resp.Diagnostics.Append(diags...)

	liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
	resp.Diagnostics.Append(diags...)

	// Build state
	state := OrganizationalResourceGroupModel{
		// For data sources, the `id` in state should be stable. Often it's the same as input.
		ShortName:      types.StringValue(resourceGroup.ID.ShortName),
		OrganizationId: types.StringValue(resourceGroup.ID.ShortName),
		DisplayName:    types.StringValue(resourceGroup.DisplayName),
		Description:    types.StringValue(resourceGroup.Description),
		Status:         types.StringValue(resourceGroup.Status),
		Icon:           types.StringValue(resourceGroup.Icon),
		MembersIds:     membersIds,
		TeamsIds:       teamsIds,
		ManagersIds:    managersIds,
		FractalsIds:    fractalsIds,
		LiveSystemsIds: liveSystemsIds,
		CreatedAt:      types.StringValue(resourceGroup.CreatedAt),
		CreatedBy:      types.StringValue(resourceGroup.CreatedBy),
		UpdatedAt:      types.StringValue(resourceGroup.UpdatedAt),
		UpdatedBy:      types.StringValue(resourceGroup.UpdatedBy),
	}

	// Write state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
