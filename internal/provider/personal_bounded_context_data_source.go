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
	_ datasource.DataSource              = &PersonalBoundedContextDataSource{}
	_ datasource.DataSourceWithConfigure = &PersonalBoundedContextDataSource{}
)

// NewPersonalBoundedContextDataSource is a helper function to simplify the provider implementation.
func NewPersonalBoundedContextDataSource() datasource.DataSource {
	return &PersonalBoundedContextDataSource{}
}

// PersonalBoundedContextDataSource is the data source implementation.
type PersonalBoundedContextDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *PersonalBoundedContextDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PersonalBoundedContextDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_bounded_context"
}

// Schema defines the schema for the data source.
func (d *PersonalBoundedContextDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

// PersonalBoundedContextModel maps bounded context schema data.
type PersonalBoundedContextModel struct {
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
func (d *PersonalBoundedContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PersonalBoundedContextModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := config.ShortName.ValueString()
	if config.ShortName.IsUnknown() || config.ShortName.IsNull() || shortName == "" {
		resp.Diagnostics.AddError(
			"Invalid Bounded Context Short Name",
			"The short_name attribute is required and must not be empty.",
		)
		return
	}

	tflog.Debug(ctx, "reading personal bounded context data source", map[string]any{"short_name": shortName})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: shortName,
	}

	resourceGroup, err := d.client.GetPersonalResourceGroup(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Personal Bounded Context",
			fmt.Sprintf("Could not read bounded context %q: %s", shortName, err),
		)
		return
	}

	if resourceGroup == nil {
		resp.Diagnostics.AddError(
			"Personal Bounded Context Not Found",
			fmt.Sprintf("No bounded context found with short_name %q.", shortName),
		)
		return
	}

	var state PersonalBoundedContextModel
	mapPersonalBoundedContextToState(ctx, resourceGroup, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
