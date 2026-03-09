package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &OrganizationalBoundedContextResource{}
	_ resource.ResourceWithConfigure = &OrganizationalBoundedContextResource{}
)

// NewOrganizationalBoundedContext is a helper function to simplify the provider implementation.
func NewOrganizationalBoundedContext() resource.Resource {
	return &OrganizationalBoundedContextResource{}
}

// OrganizationalBoundedContextResource is the resource implementation.
type OrganizationalBoundedContextResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *OrganizationalBoundedContextResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fractalCloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *fractalCloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *OrganizationalBoundedContextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_bounded_context"
}

// Schema defines the schema for the resource.
func (r *OrganizationalBoundedContextResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *OrganizationalBoundedContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationalBoundedContextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := plan.ShortName.ValueString()
	orgId := plan.OrganizationId.ValueString()
	tflog.Info(ctx, "creating organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})

	resourceGroup := fractalCloud.OrganizationalResourceGroup{
		Id: fractalCloud.ResourceGroupId{
			Type:      "Organizational",
			OwnerId:   orgId,
			ShortName: shortName,
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if err := r.client.UpsertOrganizationalResourceGroup(ctx, resourceGroup); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Organizational Bounded Context",
			fmt.Sprintf("Could not create bounded context %q in organization %q: %s", shortName, orgId, err),
		)
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "created organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})
}

// Read refreshes the Terraform state with the latest data.
func (r *OrganizationalBoundedContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationalBoundedContextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := state.ShortName.ValueString()
	orgId := state.OrganizationId.ValueString()
	tflog.Debug(ctx, "reading organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerId:   orgId,
		ShortName: shortName,
	}

	resourceGroup, err := r.client.GetOrganizationalResourceGroup(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organizational Bounded Context",
			fmt.Sprintf("Could not read bounded context %q in organization %q: %s", shortName, orgId, err),
		)
		return
	}

	if resourceGroup == nil {
		tflog.Warn(ctx, "organizational bounded context not found, removing from state", map[string]any{
			"short_name":      shortName,
			"organization_id": orgId,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapOrganizationalBoundedContextToState(ctx, resourceGroup, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *OrganizationalBoundedContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationalBoundedContextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := plan.ShortName.ValueString()
	orgId := plan.OrganizationId.ValueString()
	tflog.Info(ctx, "updating organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})

	resourceGroup := fractalCloud.OrganizationalResourceGroup{
		Id: fractalCloud.ResourceGroupId{
			Type:      "Organizational",
			OwnerId:   orgId,
			ShortName: shortName,
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if err := r.client.UpsertOrganizationalResourceGroup(ctx, resourceGroup); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Organizational Bounded Context",
			fmt.Sprintf("Could not update bounded context %q in organization %q: %s", shortName, orgId, err),
		)
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "updated organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *OrganizationalBoundedContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationalBoundedContextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := state.ShortName.ValueString()
	orgId := state.OrganizationId.ValueString()
	tflog.Info(ctx, "deleting organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerId:   orgId,
		ShortName: shortName,
	}

	if err := r.client.DeleteOrganizationalResourceGroup(ctx, resourceGroupId); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Organizational Bounded Context",
			fmt.Sprintf("Could not delete bounded context %q in organization %q: %s", shortName, orgId, err),
		)
		return
	}

	tflog.Info(ctx, "deleted organizational bounded context", map[string]any{
		"short_name":      shortName,
		"organization_id": orgId,
	})
}

// readIntoModel fetches the current state from the API and populates the model.
func (r *OrganizationalBoundedContextResource) readIntoModel(
	ctx context.Context,
	model *OrganizationalBoundedContextModel,
	diags *Diagnostics,
) {
	shortName := model.ShortName.ValueString()
	orgId := model.OrganizationId.ValueString()

	resourceGroup, err := r.client.GetOrganizationalResourceGroup(ctx, fractalCloud.ResourceGroupId{
		Type:      "Organizational",
		OwnerId:   orgId,
		ShortName: shortName,
	})
	if err != nil {
		diags.AddError(
			"Error Reading Organizational Bounded Context After Mutation",
			fmt.Sprintf("Could not read bounded context %q in organization %q after create/update: %s", shortName, orgId, err),
		)
		return
	}

	if resourceGroup == nil {
		diags.AddError(
			"Organizational Bounded Context Not Found After Mutation",
			fmt.Sprintf("Bounded context %q in organization %q was not found immediately after create/update. This is unexpected.", shortName, orgId),
		)
		return
	}

	mapOrganizationalBoundedContextToState(ctx, resourceGroup, model, diags)
}

// mapOrganizationalBoundedContextToState maps an API response to the Terraform model.
func mapOrganizationalBoundedContextToState(
	ctx context.Context,
	resourceGroup *fractalCloud.OrganizationalResourceGroup,
	model *OrganizationalBoundedContextModel,
	diags *Diagnostics,
) {
	membersIds, d := types.ListValueFrom(ctx, types.StringType, resourceGroup.MembersIds)
	diags.Append(d...)

	teamsIds, d := types.ListValueFrom(ctx, types.StringType, resourceGroup.TeamsIds)
	diags.Append(d...)

	managersIds, d := types.ListValueFrom(ctx, types.StringType, resourceGroup.ManagersIds)
	diags.Append(d...)

	fractalsIds, d := types.ListValueFrom(ctx, types.StringType, resourceGroup.FractalsIds)
	diags.Append(d...)

	liveSystemsIds, d := types.ListValueFrom(ctx, types.StringType, resourceGroup.LiveSystemsIds)
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	idTypes := map[string]attr.Type{
		"type":       types.StringType,
		"owner_id":   types.StringType,
		"short_name": types.StringType,
	}

	model.Id = types.ObjectValueMust(idTypes, map[string]attr.Value{
		"type":       types.StringValue(resourceGroup.Id.Type),
		"owner_id":   types.StringValue(resourceGroup.Id.OwnerId),
		"short_name": types.StringValue(resourceGroup.Id.ShortName),
	})
	model.ShortName = types.StringValue(resourceGroup.Id.ShortName)
	model.OrganizationId = types.StringValue(resourceGroup.Id.OwnerId)
	model.DisplayName = stringValueOrNull(resourceGroup.DisplayName, model.DisplayName)
	model.Description = stringValueOrNull(resourceGroup.Description, model.Description)
	model.Status = types.StringValue(resourceGroup.Status)
	model.Icon = stringValueOrNull(resourceGroup.Icon, model.Icon)
	model.MembersIds = membersIds
	model.TeamsIds = teamsIds
	model.ManagersIds = managersIds
	model.FractalsIds = fractalsIds
	model.LiveSystemsIds = liveSystemsIds
	model.CreatedAt = types.StringValue(resourceGroup.CreatedAt)
	model.CreatedBy = types.StringValue(resourceGroup.CreatedBy)
	model.UpdatedAt = types.StringValue(resourceGroup.UpdatedAt)
	model.UpdatedBy = types.StringValue(resourceGroup.UpdatedBy)
}
