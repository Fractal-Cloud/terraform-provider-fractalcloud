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
	_ resource.Resource              = &OrganizationalResourceGroupResource{}
	_ resource.ResourceWithConfigure = &OrganizationalResourceGroupResource{}
)

// NewOrganizationalResourceGroup is a helper function to simplify the provider implementation.
func NewOrganizationalResourceGroup() resource.Resource {
	return &OrganizationalResourceGroupResource{}
}

// OrganizationalResourceGroupResource is the resource implementation.
type OrganizationalResourceGroupResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *OrganizationalResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *OrganizationalResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_resource_group"
}

// Schema defines the schema for the resource.
func (r *OrganizationalResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"short_name": schema.StringAttribute{
				Required: true,
			},
			"organization_id": schema.StringAttribute{
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
			"members_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
			},
			"teams_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
			},
			"managers_ids": schema.ListAttribute{
				Optional:    true,
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

// Create creates the resource and sets the initial Terraform state.
func (r *OrganizationalResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan OrganizationalResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdResourceGroup, err := UpsertOrganizationalResourceGroup(ctx, resp.Diagnostics, plan, r)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Fractal Organizational Resource Group",
			"Could not create Fractal Organizational Resource Group, unexpected error: "+err.Error(),
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
func (r *OrganizationalResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state OrganizationalResourceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerID:   state.OrganizationId.ValueString(),
		ShortName: state.ShortName.ValueString(),
	}
	resourceGroup, err := r.client.GetOrganizationalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal Cloud Organizational Resource Group",
			"Could not read Fractal Cloud Organizational Resource Group with ShortName "+
				state.ShortName.ValueString()+" within the organization "+
				state.OrganizationId.ValueString()+": "+err.Error(),
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
func (r *OrganizationalResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationalResourceGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedResourceGroup, err := UpsertOrganizationalResourceGroup(ctx, resp.Diagnostics, plan, r)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating or updating Organizational Fractal Resource Group",
			"Could not create Fractal Organizational Resource Group, unexpected error: "+err.Error(),
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
func (r *OrganizationalResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state OrganizationalResourceGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerID:   state.OrganizationId.ValueString(),
		ShortName: state.ShortName.ValueString(),
	}

	// Delete existing order
	err := r.client.DeleteOrganizationalResourceGroup(resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Fractal Cloud Organizational Resource Group",
			"Could not delete Organizational Resource Group, unexpected error: "+err.Error(),
		)
		return
	}
}

func UpsertOrganizationalResourceGroup(
	ctx context.Context,
	diagnostics diag.Diagnostics,
	plan OrganizationalResourceGroupModel,
	r *OrganizationalResourceGroupResource) (*OrganizationalResourceGroupModel, error) {
	// Generate API request body from plan
	var resourceGroup = fractalCloud.OrganizationalResourceGroup{
		ID: fractalCloud.ResourceGroupId{
			Type:      "Organizational",
			OwnerID:   plan.OrganizationId.ValueString(),
			ShortName: plan.ShortName.ValueString(),
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Create new order
	err := r.client.UpsertOrganizationalResourceGroup(resourceGroup)
	if err != nil {
		return nil, err
	}

	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))

	updatedResourceGroup, err := r.client.GetOrganizationalResourceGroup(resourceGroup.ID)
	if err != nil {
		return nil, err
	}

	if updatedResourceGroup == nil {
		return nil, errors.New("organizational resource group not found after upsert")
	}

	membersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.MembersIds)
	diagnostics.Append(diags...)

	teamsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.TeamsIds)
	diagnostics.Append(diags...)

	managersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.ManagersIds)
	diagnostics.Append(diags...)

	fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
	diagnostics.Append(diags...)

	liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
	diagnostics.Append(diags...)

	return &OrganizationalResourceGroupModel{
		ShortName:      types.StringValue(resourceGroup.ID.ShortName),
		OrganizationId: types.StringValue(resourceGroup.ID.ShortName),
		DisplayName:    types.StringValue(resourceGroup.DisplayName),
		Description:    types.StringValue(resourceGroup.Description),
		Status:         types.StringValue(resourceGroup.Status),
		Icon:           types.StringValue(resourceGroup.Icon),
		MembersIds:     membersIds,
		TeamsIds:       teamsIds,
		ManagersIds:    managersIds,
		FractalsIds:    fractalsIds,
		LiveSystemsIds: liveSystemsIds,
		CreatedAt:      types.StringValue(resourceGroup.CreatedAt),
		CreatedBy:      types.StringValue(resourceGroup.CreatedBy),
		UpdatedAt:      types.StringValue(resourceGroup.UpdatedAt),
		UpdatedBy:      types.StringValue(resourceGroup.UpdatedBy),
	}, nil
}
