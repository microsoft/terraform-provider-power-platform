package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSliceToSet converts a slice of strings to a Terraform set of strings.
// If an error occurs during conversion, it returns an error instead of panicking.
func StringSliceToSet(slice []string) (types.Set, error) {
	values := make([]attr.Value, len(slice))
	for i, v := range slice {
		values[i] = types.StringValue(v)
	}
	set, diags := types.SetValue(types.StringType, values)
	if diags.HasError() {
		return types.Set{}, fmt.Errorf("failed to convert string slice to set: %v", diags)
	}
	return set, nil
}
