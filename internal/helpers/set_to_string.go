package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetToStringSlice converts a set of strings to a slice of strings.

func SetToStringSlice(set types.Set) []string {

	var result []string

	for _, v := range set.Elements() {

		if str, ok := v.(types.String); ok {

			result = append(result, str.ValueString())

		}

	}

	return result

}

// Add other helper functions here
