package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// linkAttrTypes defines the attribute types for a component link object.
var linkAttrTypes = map[string]attr.Type{
	"component_id": types.StringType,
	"settings":     types.MapType{ElemType: types.StringType},
}

// linkObjectType is the ObjectType for a single component link.
var linkObjectType = types.ObjectType{AttrTypes: linkAttrTypes}

// componentAttrTypes defines the attribute types for a component object.
// This must exactly match the nested object schema in fractal_resource.go.
var componentAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"type":                types.StringType,
	"display_name":        types.StringType,
	"description":         types.StringType,
	"version":             types.StringType,
	"is_locked":           types.BoolType,
	"recreate_on_failure": types.BoolType,
	"parameters":          types.MapType{ElemType: types.StringType},
	"dependencies_ids":    types.ListType{ElemType: types.StringType},
	"links":               types.ListType{ElemType: linkObjectType},
	"output_fields":       types.ListType{ElemType: types.StringType},
}

// componentReturn returns the standard function return type for all component functions.
func componentReturn() function.ObjectReturn {
	return function.ObjectReturn{
		AttributeTypes: componentAttrTypes,
	}
}

// buildComponent constructs a types.Object representing a blueprint component.
func buildComponent(
	id string,
	componentType string,
	displayName types.String,
	description types.String,
	version types.String,
	parameters map[string]string,
	dependenciesIds []string,
) (types.Object, *function.FuncError) {
	// Build parameters map value
	var parametersValue attr.Value
	if len(parameters) > 0 {
		elems := make(map[string]attr.Value, len(parameters))
		for k, v := range parameters {
			elems[k] = types.StringValue(v)
		}
		mv, diags := types.MapValue(types.StringType, elems)
		if diags.HasError() {
			return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build parameters map")
		}
		parametersValue = mv
	} else {
		parametersValue = types.MapNull(types.StringType)
	}

	// Build dependencies list value
	var depsValue attr.Value
	if len(dependenciesIds) > 0 {
		elems := make([]attr.Value, len(dependenciesIds))
		for i, dep := range dependenciesIds {
			elems[i] = types.StringValue(dep)
		}
		lv, diags := types.ListValue(types.StringType, elems)
		if diags.HasError() {
			return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build dependencies list")
		}
		depsValue = lv
	} else {
		depsValue = types.ListNull(types.StringType)
	}

	attrs := map[string]attr.Value{
		"id":                  types.StringValue(id),
		"type":                types.StringValue(componentType),
		"display_name":        displayName,
		"description":         description,
		"version":             version,
		"is_locked":           types.BoolNull(),
		"recreate_on_failure": types.BoolNull(),
		"parameters":          parametersValue,
		"dependencies_ids":    depsValue,
		"links":               types.ListNull(linkObjectType),
		"output_fields":       types.ListNull(types.StringType),
	}

	obj, diags := types.ObjectValue(componentAttrTypes, attrs)
	if diags.HasError() {
		return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build component object")
	}
	return obj, nil
}

// optionalString returns the types.String value if non-null, otherwise types.StringNull().
func optionalString(v types.String) types.String {
	if v.IsNull() || v.IsUnknown() {
		return types.StringNull()
	}
	return v
}
