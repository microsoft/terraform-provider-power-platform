package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSliceToSet converts a slice of strings to a Terraform set of strings.

func StringSliceToSet(slice []string) types.Set {
	values := make([]attr.Value, len(slice))
	for i, v := range slice {
		values[i] = types.StringValue(v)
	}
	set, diags := types.SetValue(types.StringType, values)
	if diags.HasError() {
		panic("failed to convert string slice to set")
	}
	return set
}
