package provider

import (
	"context"
	"fmt"
	"strconv"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
func (r *FractalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fractal"
}

// Schema defines the schema for the resource.
func (r *FractalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bounded_context_id": schema.ObjectAttribute{
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
							Computed: true,
							Default:  stringdefault.StaticString(""),
						},
						"description": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(""),
						},
						"version": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString(""),
						},
						"is_locked": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"recreate_on_failure": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"parameters": schema.MapAttribute{
							Optional:    true,
							Computed:    true,
							ElementType: basetypes.StringType{},
							Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
						},
						"dependencies_ids": schema.ListAttribute{
							Optional:    true,
							Computed:    true,
							ElementType: basetypes.StringType{},
							Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
						},
						"links": schema.ListNestedAttribute{
							Optional: true,
							Computed: true,
							Default: listdefault.StaticValue(types.ListValueMust(types.ObjectType{
								AttrTypes: linkAttrTypes,
							}, []attr.Value{})),
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
							Computed:    true,
							ElementType: basetypes.StringType{},
							Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
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

// fractalIdFromModel constructs a FractalId from the BlueprintModel.
func fractalIdFromModel(plan BlueprintModel) fractalCloud.FractalId {
	return fractalCloud.FractalId{
		ResourceGroupId: fractalCloud.ResourceGroupId{
			Type:      plan.BoundedContextId.Type.ValueString(),
			OwnerId:   plan.BoundedContextId.OwnerId.ValueString(),
			ShortName: plan.BoundedContextId.ShortName.ValueString(),
		},
		Name:    plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}
}

// fractalLogFields returns structured log fields for a fractal.
func fractalLogFields(plan BlueprintModel) map[string]any {
	return map[string]any{
		"fractal_name":    plan.Name.ValueString(),
		"fractal_version": plan.Version.ValueString(),
		"bounded_context": plan.BoundedContextId.ShortName.ValueString(),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *FractalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fractalId := fractalIdFromModel(plan)
	tflog.Info(ctx, "creating fractal", fractalLogFields(plan))

	components, err := r.extractComponents(ctx, plan, &resp.Diagnostics)
	if err != nil {
		return
	}

	tflog.Debug(ctx, "fractal create request", map[string]any{
		"fractal_id":      fractalId.ToString(),
		"component_count": len(components),
	})

	// Check if the fractal already exists — if so, update instead of create
	existing, err := r.client.GetBlueprint(ctx, fractalId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Checking Fractal Existence",
			fmt.Sprintf("Could not check if fractal %q already exists: %s", fractalId.ToString(), err),
		)
		return
	}

	if existing != nil {
		tflog.Info(ctx, "fractal already exists, updating instead", fractalLogFields(plan))
		err = r.client.UpdateBlueprint(ctx, fractalId, plan.Description.ValueString(), plan.IsPrivate.ValueBool(), components)
	} else {
		err = r.client.CreateBlueprint(ctx, fractalId, plan.Description.ValueString(), plan.IsPrivate.ValueBool(), components)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Fractal",
			fmt.Sprintf("Could not create fractal %q: %s", fractalId.ToString(), err),
		)
		return
	}

	// Only fetch the server-computed field (created_at); keep plan as source of truth
	r.readComputedFields(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "created fractal", fractalLogFields(plan))
}

// Read refreshes the Terraform state with the latest data.
func (r *FractalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fractalId := fractalIdFromModel(state)
	tflog.Debug(ctx, "reading fractal", fractalLogFields(state))

	blueprint, err := r.client.GetBlueprint(ctx, fractalId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal",
			fmt.Sprintf("Could not read fractal %q: %s", fractalId.ToString(), err),
		)
		return
	}

	if blueprint == nil {
		tflog.Warn(ctx, "fractal not found, removing from state", fractalLogFields(state))
		resp.State.RemoveResource(ctx)
		return
	}

	mapBlueprintToState(ctx, blueprint, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *FractalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BlueprintModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fractalId := fractalIdFromModel(plan)
	tflog.Info(ctx, "updating fractal", fractalLogFields(plan))

	components, err := r.extractComponents(ctx, plan, &resp.Diagnostics)
	if err != nil {
		return
	}

	tflog.Debug(ctx, "fractal update request", map[string]any{
		"fractal_id":      fractalId.ToString(),
		"component_count": len(components),
	})

	if err := r.client.UpdateBlueprint(ctx, fractalId, plan.Description.ValueString(), plan.IsPrivate.ValueBool(), components); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Fractal",
			fmt.Sprintf("Could not update fractal %q: %s", fractalId.ToString(), err),
		)
		return
	}

	r.readComputedFields(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Info(ctx, "updated fractal", fractalLogFields(plan))
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *FractalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BlueprintModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fractalId := fractalIdFromModel(state)
	tflog.Info(ctx, "deleting fractal", fractalLogFields(state))

	if err := r.client.DeleteBlueprint(ctx, fractalId); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Fractal",
			fmt.Sprintf("Could not delete fractal %q: %s", fractalId.ToString(), err),
		)
		return
	}

	tflog.Info(ctx, "deleted fractal", fractalLogFields(state))
}

// extractComponents parses the components from the plan into client model objects.
// Returns nil and adds diagnostics on error.
func (r *FractalResource) extractComponents(
	ctx context.Context,
	plan BlueprintModel,
	diags *Diagnostics,
) ([]fractalCloud.Component, error) {
	components := make([]fractalCloud.Component, 0, len(plan.Components.Elements()))
	d := plan.Components.ElementsAs(ctx, &components, false)
	diags.Append(d...)
	if d.HasError() {
		logDiags(ctx, "failed to parse fractal components", d)
		diags.AddError(
			"Invalid Fractal Components",
			"Could not parse the components configuration. Check that all component attributes have the correct types.",
		)
		return nil, fmt.Errorf("invalid components")
	}

	tflog.Debug(ctx, "parsed fractal components", map[string]any{
		"count": strconv.Itoa(len(components)),
	})
	return components, nil
}

// readComputedFields fetches only server-computed fields (created_at) from the API,
// preserving all plan values as the source of truth for non-computed attributes.
func (r *FractalResource) readComputedFields(
	ctx context.Context,
	model *BlueprintModel,
	diags *Diagnostics,
) {
	fractalId := fractalIdFromModel(*model)

	blueprint, err := r.client.GetBlueprint(ctx, fractalId)
	if err != nil {
		diags.AddError(
			"Error Reading Fractal After Mutation",
			fmt.Sprintf("Could not read fractal %q after create/update: %s", fractalId.ToString(), err),
		)
		return
	}

	if blueprint == nil {
		diags.AddError(
			"Fractal Not Found After Mutation",
			fmt.Sprintf("Fractal %q was not found immediately after create/update. This is unexpected.", fractalId.ToString()),
		)
		return
	}

	model.CreatedAt = types.StringValue(blueprint.CreatedAt)
}
