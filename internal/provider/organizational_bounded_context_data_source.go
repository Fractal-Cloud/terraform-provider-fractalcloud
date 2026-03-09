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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &OrganizationalBoundedContextDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationalBoundedContextDataSource{}
)

// NewOrganizationalBoundedContextDataSource is a helper function to simplify the provider implementation.
func NewOrganizationalBoundedContextDataSource() datasource.DataSource {
	return &OrganizationalBoundedContextDataSource{}
}

// OrganizationalBoundedContextDataSource is the data source implementation.
type OrganizationalBoundedContextDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationalBoundedContextDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *OrganizationalBoundedContextDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_bounded_context"
}

// Schema defines the schema for the data source.
func (d *OrganizationalBoundedContextDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

// OrganizationalBoundedContextModel maps bounded context schema data.
type OrganizationalBoundedContextModel struct {
	Id             types.Object `tfsdk:"id"`
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
func (d *OrganizationalBoundedContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OrganizationalBoundedContextModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := config.ShortName.ValueString()
	orgId := config.OrganizationId.ValueString()

	if config.ShortName.IsUnknown() || config.ShortName.IsNull() || shortName == "" {
		resp.Diagnostics.AddError(
			"Invalid Bounded Context Short Name",
			"The short_name attribute is required and must not be empty.",
		)
		return
	}

	tflog.Debug(ctx, "reading organizational bounded context data source", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerId:   orgId,
		ShortName: shortName,
	}

	resourceGroup, err := d.client.GetOrganizationalResourceGroup(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organizational Bounded Context",
			fmt.Sprintf("Could not read bounded context %q in organization %q: %s", shortName, orgId, err),
		)
		return
	}

	if resourceGroup == nil {
		resp.Diagnostics.AddError(
			"Organizational Bounded Context Not Found",
			fmt.Sprintf("No bounded context found with short_name %q in organization %q.", shortName, orgId),
		)
		return
	}

	var state OrganizationalBoundedContextModel
	mapOrganizationalBoundedContextToState(ctx, resourceGroup, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
