package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
func (d *OrganizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
			"bounded_contexts": schema.ListAttribute{
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

// OrganizationModel maps organization schema data.
type OrganizationModel struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Icon            types.String `tfsdk:"icon"`
	Tags            types.List   `tfsdk:"tags"`
	SocialLinks     types.List   `tfsdk:"social_links"`
	Admins          types.List   `tfsdk:"admins"`
	Members         types.List   `tfsdk:"members"`
	Teams           types.List   `tfsdk:"teams"`
	BoundedContexts types.List   `tfsdk:"bounded_contexts"`
	Status          types.String `tfsdk:"status"`
	SubscriptionId  types.String `tfsdk:"subscription_id"`
	CreatedAt       types.String `tfsdk:"created_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	UpdatedBy       types.String `tfsdk:"updated_by"`
}

// Read refreshes the Terraform state with the latest data.
func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OrganizationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := config.Id.ValueString()
	if config.Id.IsUnknown() || config.Id.IsNull() || orgId == "" {
		resp.Diagnostics.AddError(
			"Invalid Organization ID",
			"The id attribute is required and must not be empty.",
		)
		return
	}

	tflog.Debug(ctx, "reading organization data source", map[string]any{"organization_id": orgId})

	organization, err := d.client.GetOrganization(ctx, orgId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organization",
			fmt.Sprintf("Could not read organization %q: %s", orgId, err),
		)
		return
	}

	if organization == nil {
		resp.Diagnostics.AddError(
			"Organization Not Found",
			fmt.Sprintf("No organization found with id %q.", orgId),
		)
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

	boundedContexts, diags := types.ListValueFrom(ctx, types.StringType, organization.ResourceGroupsIds)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := OrganizationModel{
		Id:              types.StringValue(organization.Id),
		DisplayName:     types.StringValue(organization.DisplayName),
		Description:     types.StringValue(organization.Description),
		Icon:            types.StringValue(organization.Icon),
		Tags:            tags,
		SocialLinks:     socialLinks,
		Admins:          admins,
		Members:         members,
		Teams:           teams,
		BoundedContexts: boundedContexts,
		Status:          types.StringValue(organization.Status),
		SubscriptionId:  types.StringValue(organization.SubscriptionId),
		CreatedAt:       types.StringValue(organization.CreatedAt),
		CreatedBy:       types.StringValue(organization.CreatedBy),
		UpdatedAt:       types.StringValue(organization.UpdatedAt),
		UpdatedBy:       types.StringValue(organization.UpdatedBy),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
