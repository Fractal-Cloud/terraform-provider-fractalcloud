package provider

import (
	"fmt"

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

// portLinkAttrTypes defines the object attributes for a port-based traffic link
// used in function parameters (VM→VM, Workload→Workload).
var portLinkAttrTypes = map[string]attr.Type{
	"target_id": types.StringType,
	"from_port": types.Int64Type,
	"to_port":   types.Int64Type,
	"protocol":  types.StringType,
}

// componentLink represents a resolved link ready to be set on the component.
type componentLink struct {
	ComponentId string
	Settings    map[string]string
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
	links []componentLink,
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

	// Build links list value
	var linksValue attr.Value
	if len(links) > 0 {
		linkElems := make([]attr.Value, len(links))
		for i, link := range links {
			// Build settings map for this link
			var settingsValue attr.Value
			if len(link.Settings) > 0 {
				settingsElems := make(map[string]attr.Value, len(link.Settings))
				for k, v := range link.Settings {
					settingsElems[k] = types.StringValue(v)
				}
				sv, diags := types.MapValue(types.StringType, settingsElems)
				if diags.HasError() {
					return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build link settings map")
				}
				settingsValue = sv
			} else {
				settingsValue = types.MapNull(types.StringType)
			}

			linkObj, diags := types.ObjectValue(linkAttrTypes, map[string]attr.Value{
				"component_id": types.StringValue(link.ComponentId),
				"settings":     settingsValue,
			})
			if diags.HasError() {
				return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build link object")
			}
			linkElems[i] = linkObj
		}
		lv, diags := types.ListValue(linkObjectType, linkElems)
		if diags.HasError() {
			return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build links list")
		}
		linksValue = lv
	} else {
		linksValue = types.ListNull(linkObjectType)
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
		"links":               linksValue,
		"output_fields":       types.ListNull(types.StringType),
	}

	obj, diags := types.ObjectValue(componentAttrTypes, attrs)
	if diags.HasError() {
		return types.ObjectNull(componentAttrTypes), function.NewFuncError("failed to build component object")
	}
	return obj, nil
}

// portLinksToComponentLinks converts port-based traffic link configs to componentLinks.
func portLinksToComponentLinks(portLinks []portLinkConfig) []componentLink {
	result := make([]componentLink, len(portLinks))
	for i, pl := range portLinks {
		settings := map[string]string{
			"fromPort": fmt.Sprintf("%d", pl.FromPort.ValueInt64()),
		}

		if !pl.ToPort.IsNull() && !pl.ToPort.IsUnknown() {
			settings["toPort"] = fmt.Sprintf("%d", pl.ToPort.ValueInt64())
		} else {
			settings["toPort"] = fmt.Sprintf("%d", pl.FromPort.ValueInt64())
		}

		if !pl.Protocol.IsNull() && !pl.Protocol.IsUnknown() {
			settings["protocol"] = pl.Protocol.ValueString()
		} else {
			settings["protocol"] = "tcp"
		}

		result[i] = componentLink{
			ComponentId: pl.TargetId.ValueString(),
			Settings:    settings,
		}
	}
	return result
}

// sgMembershipLinks converts a list of security group IDs to membership componentLinks (empty settings).
func sgMembershipLinks(sgIds []string) []componentLink {
	result := make([]componentLink, len(sgIds))
	for i, sgId := range sgIds {
		result[i] = componentLink{
			ComponentId: sgId,
			Settings:    nil,
		}
	}
	return result
}

// portLinkConfig is the Go struct for a port-based traffic link function parameter.
type portLinkConfig struct {
	TargetId types.String `tfsdk:"target_id"`
	FromPort types.Int64  `tfsdk:"from_port"`
	ToPort   types.Int64  `tfsdk:"to_port"`
	Protocol types.String `tfsdk:"protocol"`
}

// optionalString returns the types.String value if non-null, otherwise types.StringNull().
func optionalString(v types.String) types.String {
	if v.IsNull() || v.IsUnknown() {
		return types.StringNull()
	}
	return v
}
