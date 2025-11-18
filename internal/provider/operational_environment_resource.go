package provider

import (
	"context"
	"fmt"

	fractal_cloud "fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &operationalEnvironment{}
	_ resource.ResourceWithConfigure = &operationalEnvironment{}
)

// NewOperationalEnvironment is a helper function to simplify the provider implementation.
func NewOperationalEnvironment() resource.Resource {
	return &operationalEnvironment{}
}

// orderResource is the resource implementation.
type operationalEnvironment struct {
	client *fractal_cloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *operationalEnvironment) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fractal_cloud.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fractal_cloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *operationalEnvironment) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operational_environment"
}

// Schema defines the schema for the resource.
func (r *operationalEnvironment) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"type":      basetypes.StringType{},
					"owner_id":  basetypes.StringType{},
					"shortname": basetypes.StringType{},
				},
			},
			"management_environment_id": schema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"type":      basetypes.StringType{},
					"owner_id":  basetypes.StringType{},
					"shortname": basetypes.StringType{},
				},
			},
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"resource_groups": schema.ListAttribute{
				Required: true,
				ElementType: basetypes.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":      basetypes.StringType{},
						"owner_id":  basetypes.StringType{},
						"shortname": basetypes.StringType{},
					},
				},
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				ElementType: basetypes.ObjectType{},
			},
			"agents": schema.SetAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
			},
			"default_cicd_profile_short_name": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
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

// Create creates the resource and sets the initial Terraform state.
func (r *operationalEnvironment) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *operationalEnvironment) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *operationalEnvironment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *operationalEnvironment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
