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
	_ datasource.DataSource              = &ResourceGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &ResourceGroupsDataSource{}
)

// NewResourceGroupDataSource is a helper function to simplify the provider implementation.
func NewResourceGroupDataSource() datasource.DataSource {
	return &ResourceGroupsDataSource{}
}

// ResourceGroupsDataSource is the data source implementation.
type ResourceGroupsDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *ResourceGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ResourceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

// Schema defines the schema for the data source.
func (d *ResourceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"type":      basetypes.StringType{},
					"owner_id":  basetypes.StringType{},
					"shortname": basetypes.StringType{},
				},
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

// ResourceGroupModel maps coffees schema data.
type ResourceGroupModel struct {
	ID          ResourceGroupIdModel `tfsdk:"id"`
	DisplayName types.String         `tfsdk:"display_name"`
	Description types.String         `tfsdk:"description"`
	Icon        types.String         `tfsdk:"icon"`
	CreatedAt   types.String         `tfsdk:"created_at"`
	CreatedBy   types.String         `tfsdk:"created_by"`
	UpdatedAt   types.String         `tfsdk:"updated_at"`
	UpdatedBy   types.String         `tfsdk:"updated_by"`
}

// ResourceGroupIdModel maps coffee ingredients data
type ResourceGroupIdModel struct {
	Type      types.String `tfsdk:"type"`
	OwnerId   types.String `tfsdk:"owner_id"`
	ShortName types.String `tfsdk:"shortname"`
}

func (id *ResourceGroupIdModel) ValueString() string {
	return id.Type.ValueString() + "/" + id.OwnerId.ValueString() + "/" + id.ShortName.ValueString()
}

// Read refreshes the Terraform state with the latest data.
func (d *ResourceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ResourceGroupModel

	// Read user configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required input isn't unknown/null
	resp.Diagnostics.Append(validateRequiredString(config.ID, "id")...)
	if resp.Diagnostics.HasError() {
		return
	}
	var resourceGroupId = fractalCloud.ResourceGroupId{
		Type:      config.ID.Type.ValueString(),
		OwnerId:   config.ID.OwnerId.ValueString(),
		ShortName: config.ID.ShortName.ValueString(),
	}

	resourceGroup, err := d.client.GetResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not read Fractal Cloud Resource Group ID "+config.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Build state
	state := ResourceGroupModel{
		// For data sources, the `id` in state should be stable. Often it's the same as input.
		ID: ResourceGroupIdModel{
			Type:      types.StringValue(resourceGroup.ID.Type),
			OwnerId:   types.StringValue(resourceGroup.ID.OwnerId),
			ShortName: types.StringValue(resourceGroup.ID.ShortName),
		},
		DisplayName: types.StringValue(resourceGroup.DisplayName),
		Description: types.StringValue(resourceGroup.Description),
		Icon:        types.StringValue(resourceGroup.Icon),
		CreatedAt:   types.StringValue(resourceGroup.CreatedAt),
		CreatedBy:   types.StringValue(resourceGroup.CreatedBy),
		UpdatedAt:   types.StringValue(resourceGroup.UpdatedAt),
		UpdatedBy:   types.StringValue(resourceGroup.UpdatedBy),
	}

	// Write state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func validateRequiredString(v ResourceGroupIdModel, attr string) diag.Diagnostics {
	var diags diag.Diagnostics

	if v.Type.IsUnknown() || v.Type.IsNull() {
		diags.AddError(
			"Unknown required value",
			fmt.Sprintf("Attribute %q.type is required, but its value is unknown, null or empty.", attr),
		)
		return diags
	}

	if v.OwnerId.IsUnknown() || v.OwnerId.IsNull() {
		diags.AddError(
			"Unknown required value",
			fmt.Sprintf("Attribute %q.owner_id is required, but its value is unknown, null or empty.", attr),
		)
		return diags
	}

	if v.ShortName.IsUnknown() || v.ShortName.IsNull() {
		diags.AddError(
			"Unknown required value",
			fmt.Sprintf("Attribute %q.shortname is required, but its value is unknown, null or empty.", attr),
		)
		return diags
	}

	return diags
}
