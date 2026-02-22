package provider

import (
	"context"
	"errors"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &FractalResource{}
	_ resource.ResourceWithConfigure = &FractalResource{}
)

// NewFractal is a helper function to simplify the provider implementation.
func NewFractal() resource.Resource {
	return &FractalResource{}
}

// FractalResource is the resource implementation.
type FractalResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *FractalResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *FractalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fractal"
}

// Schema defines the schema for the resource.
func (r *FractalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"is_private": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"components": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"display_name": schema.StringAttribute{
							Optional: true,
						},
						"description": schema.StringAttribute{
							Optional: true,
						},
						"version": schema.StringAttribute{
							Optional: true,
						},
						"is_locked": schema.StringAttribute{
							Optional: true,
						},
						"recreate_on_failure": schema.StringAttribute{
							Optional: true,
						},
						"parameters": schema.MapAttribute{
							Optional:    true,
							ElementType: basetypes.StringType{},
						},
						"dependencies_ids": schema.ListAttribute{
							Optional:    true,
							ElementType: basetypes.StringType{},
						},
						"links": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"component_id": schema.StringAttribute{
										Required: true,
									},
									"settings": schema.MapAttribute{
										Optional:    true,
										ElementType: basetypes.StringType{},
									},
								},
							},
						},
						"output_fields": schema.ListAttribute{
							Optional:    true,
							ElementType: basetypes.StringType{},
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

// Create creates the resource and sets the initial Terraform state.
func (r *FractalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Personal Resource Group")
	// Retrieve values from plan
	var plan BlueprintModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdFractal, err := createOrUpdateFractal(ctx, diags, plan, r.client, r.client.CreateBlueprint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Fractal",
			"Could not create Fractal, unexpected error: "+err.Error(),
		)
	}

	if createdFractal != nil {
		plan.ResourceGroupId = createdFractal.ResourceGroupId
		plan.Name = createdFractal.Name
		plan.Version = createdFractal.Version
		plan.Description = createdFractal.Description
		plan.IsPrivate = createdFractal.IsPrivate
		plan.Components = createdFractal.Components
		plan.CreatedAt = createdFractal.CreatedAt
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *FractalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state BlueprintModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, done := GetFractalModel(ctx, state, r.client, resp.Diagnostics)
	if done {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *FractalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Creating Personal Resource Group")
	// Retrieve values from plan
	var plan BlueprintModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdFractal, err := createOrUpdateFractal(ctx, diags, plan, r.client, r.client.UpdateBlueprint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Fractal",
			"Could not create Fractal, unexpected error: "+err.Error(),
		)
	}

	if createdFractal != nil {
		plan.ResourceGroupId = createdFractal.ResourceGroupId
		plan.Name = createdFractal.Name
		plan.Version = createdFractal.Version
		plan.Description = createdFractal.Description
		plan.IsPrivate = createdFractal.IsPrivate
		plan.Components = createdFractal.Components
		plan.CreatedAt = createdFractal.CreatedAt
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *FractalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state PersonalResourceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: state.ShortName.ValueString(),
	}

	// Delete existing order
	err := r.client.DeletePersonalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Fractal Cloud Resource Group",
			"Could not delete Resource Group, unexpected error: "+err.Error(),
		)
		return
	}
}

func createOrUpdateFractal(
	ctx context.Context,
	diagnostics diag.Diagnostics,
	plan BlueprintModel,
	client *fractalCloud.Client,
	createOrUpdate func(fractalCloud.FractalId, string, bool, []fractalCloud.Component) error) (*BlueprintModel, error) {

	// Generate API request body from plan
	var fractalId = fractalCloud.FractalId{
		ResourceGroupId: fractalCloud.ResourceGroupId{
			Type:      plan.ResourceGroupId.Type.ValueString(),
			OwnerId:   plan.ResourceGroupId.OwnerId.ValueString(),
			ShortName: plan.ResourceGroupId.ShortName.ValueString(),
		},
		Name:    plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}

	var components []fractalCloud.Component
	diags := plan.Components.ElementsAs(ctx, &components, false)
	diagnostics.Append(diags...)

	// Create or update Blueprint
	err := createOrUpdate(
		fractalId,
		plan.Description.ValueString(),
		plan.IsPrivate.ValueBool(),
		components)
	if err != nil {
		return nil, err
	}

	result, done := GetFractalModel(ctx, plan, client, diagnostics)
	if done {
		return nil, errors.New("get Fractal Model failed")
	}

	return &result, nil
}
