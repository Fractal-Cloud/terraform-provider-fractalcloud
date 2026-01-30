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
	_ datasource.DataSource              = &OrganizationDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationDataSource{}
)

// NewOrganizationDataSource is a helper function to simplify the provider implementation.
func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// OrganizationDataSource is the data source implementation.
type OrganizationDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the data source.
func (d *OrganizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"icon": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"social_links": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"admins": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"members": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"teams": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"resource_groups": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.StringType{},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"subscription_id": schema.StringAttribute{
				Computed: true,
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

// OrganizationModel maps resource group schema data.
type OrganizationModel struct {
	ID             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Icon           types.String `tfsdk:"icon"`
	Tags           types.List   `tfsdk:"tags"`
	SocialLinks    types.List   `tfsdk:"social_links"`
	Admins         types.List   `tfsdk:"admins"`
	Members        types.List   `tfsdk:"members"`
	Teams          types.List   `tfsdk:"teams"`
	ResourceGroups types.List   `tfsdk:"resource_groups"`
	Status         types.String `tfsdk:"status"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
	UpdatedBy      types.String `tfsdk:"updated_by"`
}

// Read refreshes the Terraform state with the latest data.
func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OrganizationModel

	// Read user configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required input isn't unknown/null
	if config.ID.IsUnknown() || config.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Unknown required value",
			fmt.Sprintf("ID is required, but its value is unknown, null or empty."),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	organization, err := d.client.GetOrganization(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Organization",
			"Could not read Fractal Organization with ID "+config.ID.ValueString()+": "+err.Error())
		return
	}

	if organization == nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Organization",
			"Could not find Fractal Organization with ID "+config.ID.ValueString())
		return
	}

	tags, diags := types.ListValueFrom(ctx, types.StringType, organization.Tags)
	resp.Diagnostics.Append(diags...)

	socialLinks, diags := types.ListValueFrom(ctx, types.StringType, organization.SocialLinks)
	resp.Diagnostics.Append(diags...)

	admins, diags := types.ListValueFrom(ctx, types.StringType, organization.AdminsIds)
	resp.Diagnostics.Append(diags...)

	members, diags := types.ListValueFrom(ctx, types.StringType, organization.MembersIds)
	resp.Diagnostics.Append(diags...)

	teams, diags := types.ListValueFrom(ctx, types.StringType, organization.TeamsIds)
	resp.Diagnostics.Append(diags...)

	resourceGroups, diags := types.ListValueFrom(ctx, types.StringType, organization.ResourceGroupsIds)
	resp.Diagnostics.Append(diags...)

	// Build state
	state := OrganizationModel{
		// For data sources, the `id` in state should be stable. Often it's the same as input.
		ID:             types.StringValue(organization.ID),
		DisplayName:    types.StringValue(organization.DisplayName),
		Description:    types.StringValue(organization.Description),
		Icon:           types.StringValue(organization.Icon),
		Tags:           tags,
		SocialLinks:    socialLinks,
		Admins:         admins,
		Members:        members,
		Teams:          teams,
		ResourceGroups: resourceGroups,
		Status:         types.StringValue(organization.Status),
		SubscriptionID: types.StringValue(organization.SubscriptionId),
		CreatedAt:      types.StringValue(organization.CreatedAt),
		CreatedBy:      types.StringValue(organization.CreatedBy),
		UpdatedAt:      types.StringValue(organization.UpdatedAt),
		UpdatedBy:      types.StringValue(organization.UpdatedBy),
	}

	// Write state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
