package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Diagnostics is a type alias to avoid import stutter in resource files.
type Diagnostics = diag.Diagnostics

// stringValueOrNull returns a types.StringValue for non-empty strings,
// or preserves the existing model value if the API returned empty.
// This prevents Terraform from seeing a null-to-"" diff on optional fields.
func stringValueOrNull(apiValue string, currentModelValue types.String) types.String {
	if apiValue == "" && currentModelValue.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(apiValue)
}

// stringPointerToTFValue converts a *string to a types.String.
// Returns empty string if the pointer is nil, so Terraform never sees ""→null drift.
func stringPointerToTFValue(v *string) types.String {
	if v == nil {
		return types.StringValue("")
	}
	return types.StringValue(*v)
}

// boolPointerToTFValue converts a *bool to a types.Bool.
// Returns false if the pointer is nil, so Terraform never sees false→null drift.
func boolPointerToTFValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolValue(false)
	}
	return types.BoolValue(*v)
}

func logDiags(ctx context.Context, prefix string, diags diag.Diagnostics) {
	if !diags.HasError() {
		return
	}

	tflog.Error(ctx, prefix+" (raw)", map[string]any{
		"diags_go": fmt.Sprintf("%#v", diags),
		"diags_s":  fmt.Sprintf("%v", diags),
	})

	for i, d := range diags {
		fields := map[string]any{
			"i":        i,
			"severity": d.Severity().String(),
			"summary":  d.Summary(),
			"detail":   d.Detail(),
		}
		tflog.Error(ctx, prefix, fields)
	}
}
