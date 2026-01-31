package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &PersonalResourceGroupResource{}
	_ resource.ResourceWithConfigure = &PersonalResourceGroupResource{}
)

// NewPersonalResourceGroup is a helper function to simplify the provider implementation.
func NewPersonalResourceGroup() resource.Resource {
	return &PersonalResourceGroupResource{}
}

// PersonalResourceGroupResource is the resource implementation.
type PersonalResourceGroupResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *PersonalResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PersonalResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_resource_group"
}

// Schema defines the schema for the resource.
func (r *PersonalResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"short_name": schema.StringAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"icon": schema.StringAttribute{
				Optional: true,
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

// Create creates the resource and sets the initial Terraform state.
func (r *PersonalResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan PersonalResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdResourceGroup, err := UpsertResourceGroup(ctx, diags, plan, r)
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
		plan.Status = createdResourceGroup.Status
		plan.Icon = createdResourceGroup.Icon
		plan.FractalsIds = createdResourceGroup.FractalsIds
		plan.LiveSystemsIds = createdResourceGroup.LiveSystemsIds
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *PersonalResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
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
	resourceGroup, err := r.client.GetPersonalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Resource Group",
			"Could not read Fractal Cloud Resource Group ID "+state.ShortName.ValueString()+": "+err.Error(),
		)
		return
	}

	if resourceGroup != nil {
		// Overwrite state
		fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
		resp.Diagnostics.Append(diags...)

		liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
		resp.Diagnostics.Append(diags...)

		state.DisplayName = types.StringValue(resourceGroup.DisplayName)
		state.Description = types.StringValue(resourceGroup.Description)
		state.Status = types.StringValue(resourceGroup.Status)
		state.Icon = types.StringValue(resourceGroup.Icon)
		state.FractalsIds = fractalsIds
		state.LiveSystemsIds = liveSystemsIds
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
func (r *PersonalResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PersonalResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedResourceGroup, err := UpsertResourceGroup(ctx, diags, plan, r)
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
		plan.Status = updatedResourceGroup.Status
		plan.Icon = updatedResourceGroup.Icon
		plan.FractalsIds = updatedResourceGroup.FractalsIds
		plan.LiveSystemsIds = updatedResourceGroup.LiveSystemsIds
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PersonalResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func UpsertResourceGroup(
	ctx context.Context,
	diagnostics diag.Diagnostics,
	plan PersonalResourceGroupModel,
	r *PersonalResourceGroupResource) (*PersonalResourceGroupModel, error) {
	// Generate API request body from plan
	var resourceGroup = fractalCloud.PersonalResourceGroup{
		ID: fractalCloud.ResourceGroupId{
			Type:      "Personal",
			ShortName: plan.ShortName.ValueString(),
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Create new order
	err := r.client.UpsertPersonalResourceGroup(resourceGroup)
	if err != nil {
		return nil, err
	}

	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))

	updatedResourceGroup, err := r.client.GetPersonalResourceGroup(resourceGroup.ID)
	if err != nil {
		return nil, err
	}

	if updatedResourceGroup == nil {
		return nil, errors.New("resource group not found after upsert")
	}

	fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.FractalsIds)
	diagnostics.Append(diags...)

	liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.LiveSystemsIds)
	diagnostics.Append(diags...)

	var result = &PersonalResourceGroupModel{
		ShortName:      plan.ShortName,
		DisplayName:    types.StringValue(updatedResourceGroup.DisplayName),
		Description:    types.StringValue(updatedResourceGroup.Description),
		Status:         types.StringValue(updatedResourceGroup.Status),
		LiveSystemsIds: liveSystemsIds,
		FractalsIds:    fractalsIds,
		CreatedAt:      types.StringValue(updatedResourceGroup.CreatedAt),
		UpdatedAt:      types.StringValue(updatedResourceGroup.UpdatedAt),
	}

	if len(updatedResourceGroup.Icon) > 0 {
		result.Icon = types.StringValue(updatedResourceGroup.Icon)
	} else {
		result.Icon = plan.Icon
	}

	return result, nil
}
