package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &ResourceGroupResource{}
	_ resource.ResourceWithConfigure = &ResourceGroupResource{}
)

// NewResourceGroup is a helper function to simplify the provider implementation.
func NewResourceGroup() resource.Resource {
	return &ResourceGroupResource{}
}

// ResourceGroupResource is the resource implementation.
type ResourceGroupResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *ResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

// Schema defines the schema for the resource.
func (r *ResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
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

// Create creates the resource and sets the initial Terraform state.
func (r *ResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdResourceGroup, err := UpsertResourceGroup(plan, r)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Fractal Resource Group",
			"Could not create Fractal Resource Group, unexpected error: "+err.Error(),
		)
	}

	if createdResourceGroup != nil {
		plan.CreatedAt = createdResourceGroup.CreatedAt
		plan.Description = createdResourceGroup.Description
		plan.UpdatedAt = createdResourceGroup.UpdatedAt
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ResourceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      state.ID.Type.ValueString(),
		OwnerId:   state.ID.OwnerId.ValueString(),
		ShortName: state.ID.ShortName.ValueString(),
	}
	resourceGroup, err := r.client.GetResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not read Fractal Cloud Resource Group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if resourceGroup != nil {
		// Overwrite state
		state.DisplayName = types.StringValue(resourceGroup.DisplayName)
		state.Description = types.StringValue(resourceGroup.Description)
		state.CreatedAt = types.StringValue(resourceGroup.CreatedAt)
		state.UpdatedAt = types.StringValue(resourceGroup.UpdatedAt)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedResourceGroup, err := UpsertResourceGroup(plan, r)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Fractal Resource Group",
			"Could not create Fractal Resource Group, unexpected error: "+err.Error(),
		)
	}

	if updatedResourceGroup != nil {
		plan.CreatedAt = updatedResourceGroup.CreatedAt
		plan.Description = updatedResourceGroup.Description
		plan.UpdatedAt = updatedResourceGroup.UpdatedAt
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ResourceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      state.ID.Type.ValueString(),
		OwnerId:   state.ID.OwnerId.ValueString(),
		ShortName: state.ID.ShortName.ValueString(),
	}

	// Delete existing order
	err := r.client.DeleteResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Fractal Cloud Resource Group",
			"Could not delete Resource Group, unexpected error: "+err.Error(),
		)
		return
	}
}

func UpsertResourceGroup(plan ResourceGroupModel, r *ResourceGroupResource) (*ResourceGroupModel, error) {
	// Generate API request body from plan
	var resourceGroup = fractalCloud.ResourceGroup{
		ID: fractalCloud.ResourceGroupId{
			Type:      plan.ID.Type.ValueString(),
			OwnerId:   plan.ID.OwnerId.ValueString(),
			ShortName: plan.ID.ShortName.ValueString(),
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Create new order
	err := r.client.UpsertResourceGroup(resourceGroup)
	if err != nil {
		return nil, err
	}

	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))

	updatedResourceGroup, err := r.client.GetResourceGroup(resourceGroup.ID)
	if err != nil {
		return nil, err
	}

	if updatedResourceGroup == nil {
		return nil, errors.New("resource group not found after upsert")
	}

	return &ResourceGroupModel{
		ID:          plan.ID,
		DisplayName: types.StringValue(updatedResourceGroup.DisplayName),
		Description: types.StringValue(updatedResourceGroup.Description),
		CreatedAt:   types.StringValue(updatedResourceGroup.CreatedAt),
		UpdatedAt:   types.StringValue(updatedResourceGroup.UpdatedAt),
	}, nil
}
