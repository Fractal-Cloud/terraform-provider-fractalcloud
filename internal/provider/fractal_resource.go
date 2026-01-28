package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &Fractal{}
	_ resource.ResourceWithConfigure = &Fractal{}
)

// NewFractal is a helper function to simplify the provider implementation.
func NewFractal() resource.Resource {
	return &Fractal{}
}

// Fractal orderResource is the resource implementation.
type Fractal struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *Fractal) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fractalCloud.Client)

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
func (r *Fractal) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fractal"
}

// Schema defines the schema for the resource.
func (r *Fractal) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.ObjectAttribute{
				Computed: true,
				AttributeTypes: map[string]attr.Type{
					"resource_group_id": basetypes.ObjectType{
						AttrTypes: map[string]attr.Type{
							"type":      basetypes.StringType{},
							"owner_id":  basetypes.StringType{},
							"shortname": basetypes.StringType{},
						},
					},
					"name":    basetypes.StringType{},
					"version": basetypes.StringType{},
				},
			},
			"resource_group_id": schema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"type":      basetypes.StringType{},
					"owner_id":  basetypes.StringType{},
					"shortname": basetypes.StringType{},
				},
			},
			"version": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"blueprint_components": schema.ListAttribute{
				Required: true,
				ElementType: basetypes.ObjectType{
					AttrTypes: map[string]attr.Type{
						"display_name": basetypes.StringType{},
						"type":         basetypes.StringType{},
						"id":           basetypes.StringType{},
						"version":      basetypes.StringType{},
						"locked":       basetypes.BoolType{},
						"parameters": basetypes.MapType{
							ElemType: basetypes.ObjectType{},
						},
						"links": basetypes.ListType{
							ElemType: basetypes.ObjectType{
								AttrTypes: map[string]attr.Type{
									"component_id": basetypes.StringType{},
									"settings": basetypes.MapType{
										ElemType: basetypes.ObjectType{},
									},
								},
							},
						},
					},
				},
			},
			"interface_operations": schema.ListAttribute{
				Required: true,
				ElementType: basetypes.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": basetypes.StringType{},
						"parameters": basetypes.SetType{
							ElemType: basetypes.StringType{},
						},
						"instructions": basetypes.ListType{
							ElemType: basetypes.ObjectType{
								AttrTypes: map[string]attr.Type{
									"component_id": basetypes.StringType{},
									"operation":    basetypes.StringType{},
									"input_parameters": basetypes.SetType{
										ElemType: basetypes.StringType{},
									},
								},
							},
						},
					},
				},
			},
			"private": schema.BoolAttribute{
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
func (r *Fractal) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *Fractal) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *Fractal) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *Fractal) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
