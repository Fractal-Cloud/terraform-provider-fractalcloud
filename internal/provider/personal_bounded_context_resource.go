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
	_ resource.Resource              = &PersonalBoundedContextResource{}
	_ resource.ResourceWithConfigure = &PersonalBoundedContextResource{}
)

// NewPersonalBoundedContext is a helper function to simplify the provider implementation.
func NewPersonalBoundedContext() resource.Resource {
	return &PersonalBoundedContextResource{}
}

// PersonalBoundedContextResource is the resource implementation.
type PersonalBoundedContextResource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the resource.
func (r *PersonalBoundedContextResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PersonalBoundedContextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_bounded_context"
}

// Schema defines the schema for the resource.
func (r *PersonalBoundedContextResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *PersonalBoundedContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PersonalBoundedContextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := plan.ShortName.ValueString()
	tflog.Info(ctx, "creating personal bounded context", map[string]any{"short_name": shortName})

	resourceGroup := fractalCloud.PersonalResourceGroup{
		Id: fractalCloud.ResourceGroupId{
			Type:      "Personal",
			ShortName: shortName,
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if err := r.client.UpsertPersonalResourceGroup(ctx, resourceGroup); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Personal Bounded Context",
			fmt.Sprintf("Could not create bounded context %q: %s", shortName, err),
		)
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "created personal bounded context", map[string]any{"short_name": shortName})
}

// Read refreshes the Terraform state with the latest data.
func (r *PersonalBoundedContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PersonalBoundedContextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := state.ShortName.ValueString()
	tflog.Debug(ctx, "reading personal bounded context", map[string]any{"short_name": shortName})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: shortName,
	}

	resourceGroup, err := r.client.GetPersonalResourceGroup(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Personal Bounded Context",
			fmt.Sprintf("Could not read bounded context %q: %s", shortName, err),
		)
		return
	}

	if resourceGroup == nil {
		tflog.Warn(ctx, "personal bounded context not found, removing from state", map[string]any{"short_name": shortName})
		resp.State.RemoveResource(ctx)
		return
	}

	mapPersonalBoundedContextToState(ctx, resourceGroup, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PersonalBoundedContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PersonalBoundedContextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := plan.ShortName.ValueString()
	tflog.Info(ctx, "updating personal bounded context", map[string]any{"short_name": shortName})

	resourceGroup := fractalCloud.PersonalResourceGroup{
		Id: fractalCloud.ResourceGroupId{
			Type:      "Personal",
			ShortName: shortName,
		},
		DisplayName: plan.DisplayName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if err := r.client.UpsertPersonalResourceGroup(ctx, resourceGroup); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Personal Bounded Context",
			fmt.Sprintf("Could not update bounded context %q: %s", shortName, err),
		)
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "updated personal bounded context", map[string]any{"short_name": shortName})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PersonalBoundedContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PersonalBoundedContextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shortName := state.ShortName.ValueString()
	tflog.Info(ctx, "deleting personal bounded context", map[string]any{"short_name": shortName})

	resourceGroupId := fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: shortName,
	}

	if err := r.client.DeletePersonalResourceGroup(ctx, resourceGroupId); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Personal Bounded Context",
			fmt.Sprintf("Could not delete bounded context %q: %s", shortName, err),
		)
		return
	}

	tflog.Info(ctx, "deleted personal bounded context", map[string]any{"short_name": shortName})
}

// readIntoModel fetches the current state from the API and populates the model.
// Used after Create and Update to read back server-computed fields.
func (r *PersonalBoundedContextResource) readIntoModel(
	ctx context.Context,
	model *PersonalBoundedContextModel,
	diags *Diagnostics,
) {
	shortName := model.ShortName.ValueString()

	resourceGroup, err := r.client.GetPersonalResourceGroup(ctx, fractalCloud.ResourceGroupId{
		Type:      "Personal",
		ShortName: shortName,
	})
	if err != nil {
		diags.AddError(
			"Error Reading Personal Bounded Context After Mutation",
			fmt.Sprintf("Could not read bounded context %q after create/update: %s", shortName, err),
		)
		return
	}

	if resourceGroup == nil {
		diags.AddError(
			"Personal Bounded Context Not Found After Mutation",
			fmt.Sprintf("Bounded context %q was not found immediately after create/update. This is unexpected.", shortName),
		)
		return
	}

	mapPersonalBoundedContextToState(ctx, resourceGroup, model, diags)
}

// mapPersonalBoundedContextToState maps an API response to the Terraform model.
func mapPersonalBoundedContextToState(
	ctx context.Context,
	resourceGroup *fractalCloud.PersonalResourceGroup,
	model *PersonalBoundedContextModel,
	diags *Diagnostics,
) {
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
	model.DisplayName = stringValueOrNull(resourceGroup.DisplayName, model.DisplayName)
	model.Description = stringValueOrNull(resourceGroup.Description, model.Description)
	model.Status = types.StringValue(resourceGroup.Status)
	model.Icon = stringValueOrNull(resourceGroup.Icon, model.Icon)
	model.FractalsIds = fractalsIds
	model.LiveSystemsIds = liveSystemsIds
	model.CreatedAt = types.StringValue(resourceGroup.CreatedAt)
	model.UpdatedAt = types.StringValue(resourceGroup.UpdatedAt)
}
