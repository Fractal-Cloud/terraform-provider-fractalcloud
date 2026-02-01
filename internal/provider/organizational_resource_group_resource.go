package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
		plan.Id = createdResourceGroup.Id
		plan.DisplayName = createdResourceGroup.DisplayName
		plan.Description = createdResourceGroup.Description
		plan.Status = createdResourceGroup.Status
		plan.Icon = createdResourceGroup.Icon
		plan.MembersIds = createdResourceGroup.MembersIds
		plan.ManagersIds = createdResourceGroup.ManagersIds
		plan.TeamsIds = createdResourceGroup.TeamsIds
		plan.FractalsIds = createdResourceGroup.FractalsIds
		plan.LiveSystemsIds = createdResourceGroup.LiveSystemsIds
		plan.CreatedAt = createdResourceGroup.CreatedAt
		plan.CreatedBy = createdResourceGroup.CreatedBy
		plan.UpdatedAt = createdResourceGroup.UpdatedAt
		plan.UpdatedBy = createdResourceGroup.UpdatedBy
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
		OwnerId:   state.OrganizationId.ValueString(),
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
		membersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.MembersIds)
		resp.Diagnostics.Append(diags...)

		teamsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.TeamsIds)
		resp.Diagnostics.Append(diags...)

		managersIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.ManagersIds)
		resp.Diagnostics.Append(diags...)

		fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
		resp.Diagnostics.Append(diags...)

		liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
		resp.Diagnostics.Append(diags...)

		idTypes := map[string]attr.Type{
			"type":       types.StringType,
			"owner_id":   types.StringType,
			"short_name": types.StringType,
		}

		// Overwrite state
		state.Id = types.ObjectValueMust(idTypes, map[string]attr.Value{
			"type":       types.StringValue(resourceGroup.Id.Type),
			"owner_id":   types.StringValue(resourceGroup.Id.OwnerId),
			"short_name": types.StringValue(resourceGroup.Id.ShortName),
		})
		state.DisplayName = types.StringValue(resourceGroup.DisplayName)
		state.Description = types.StringValue(resourceGroup.Description)
		state.Status = types.StringValue(resourceGroup.Status)
		state.Icon = types.StringValue(resourceGroup.Icon)
		state.MembersIds = membersIds
		state.TeamsIds = teamsIds
		state.ManagersIds = managersIds
		state.FractalsIds = fractalsIds
		state.LiveSystemsIds = liveSystemsIds
		state.CreatedAt = types.StringValue(resourceGroup.CreatedAt)
		state.CreatedBy = types.StringValue(resourceGroup.CreatedBy)
		state.UpdatedAt = types.StringValue(resourceGroup.UpdatedAt)
		state.UpdatedBy = types.StringValue(resourceGroup.UpdatedBy)
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
		plan.Id = updatedResourceGroup.Id
		plan.DisplayName = updatedResourceGroup.DisplayName
		plan.Description = updatedResourceGroup.Description
		plan.Status = updatedResourceGroup.Status
		plan.Icon = updatedResourceGroup.Icon
		plan.MembersIds = updatedResourceGroup.MembersIds
		plan.ManagersIds = updatedResourceGroup.ManagersIds
		plan.TeamsIds = updatedResourceGroup.TeamsIds
		plan.FractalsIds = updatedResourceGroup.FractalsIds
		plan.LiveSystemsIds = updatedResourceGroup.LiveSystemsIds
		plan.CreatedAt = updatedResourceGroup.CreatedAt
		plan.CreatedBy = updatedResourceGroup.CreatedBy
		plan.UpdatedAt = updatedResourceGroup.UpdatedAt
		plan.UpdatedBy = updatedResourceGroup.UpdatedBy
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
		OwnerId:   state.OrganizationId.ValueString(),
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
		Id: fractalCloud.ResourceGroupId{
			Type:      "Organizational",
			OwnerId:   plan.OrganizationId.ValueString(),
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

	updatedResourceGroup, err := r.client.GetOrganizationalResourceGroup(resourceGroup.Id)
	if err != nil {
		return nil, err
	}

	if updatedResourceGroup == nil {
		return nil, errors.New("organizational resource group not found after upsert")
	}

	membersIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.MembersIds)
	diagnostics.Append(diags...)

	teamsIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.TeamsIds)
	diagnostics.Append(diags...)

	managersIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.ManagersIds)
	diagnostics.Append(diags...)

	fractalsIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.FractalsIds)
	diagnostics.Append(diags...)

	liveSystemsIds, diags := types.ListValueFrom(ctx, types.StringType, updatedResourceGroup.LiveSystemsIds)
	diagnostics.Append(diags...)

	idTypes := map[string]attr.Type{
		"type":       types.StringType,
		"owner_id":   types.StringType,
		"short_name": types.StringType,
	}

	var result = &OrganizationalResourceGroupModel{
		Id: types.ObjectValueMust(idTypes, map[string]attr.Value{
			"type":       types.StringValue(resourceGroup.Id.Type),
			"owner_id":   types.StringValue(resourceGroup.Id.OwnerId),
			"short_name": types.StringValue(resourceGroup.Id.ShortName),
		}),
		ShortName:      types.StringValue(updatedResourceGroup.Id.ShortName),
		OrganizationId: types.StringValue(updatedResourceGroup.Id.ShortName),
		DisplayName:    types.StringValue(updatedResourceGroup.DisplayName),
		Description:    types.StringValue(updatedResourceGroup.Description),
		Status:         types.StringValue(updatedResourceGroup.Status),
		FractalsIds:    fractalsIds,
		LiveSystemsIds: liveSystemsIds,
		CreatedAt:      types.StringValue(updatedResourceGroup.CreatedAt),
		CreatedBy:      types.StringValue(updatedResourceGroup.CreatedBy),
		UpdatedAt:      types.StringValue(updatedResourceGroup.UpdatedAt),
		UpdatedBy:      types.StringValue(updatedResourceGroup.UpdatedBy),
	}

	if len(updatedResourceGroup.Icon) > 0 {
		result.Icon = types.StringValue(updatedResourceGroup.Icon)
	} else {
		result.Icon = plan.Icon
	}

	if len(updatedResourceGroup.MembersIds) > 0 {
		result.MembersIds = membersIds
	} else {
		result.MembersIds = plan.MembersIds
	}

	if len(updatedResourceGroup.TeamsIds) > 0 {
		result.TeamsIds = teamsIds
	} else {
		result.TeamsIds = plan.TeamsIds
	}

	if len(updatedResourceGroup.ManagersIds) > 0 {
		result.ManagersIds = managersIds
	} else {
		result.ManagersIds = plan.ManagersIds
	}

	return result, nil
}
