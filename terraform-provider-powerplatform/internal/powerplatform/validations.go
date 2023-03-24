package powerplatform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

//Code from:
//https://github.com/hashicorp/terraform-provider-azurerm/blob/91c25f34b2856bfd0b21978359030fd8b196c8cf/internal/tf/validation/pluginsdk.go#L174

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) ([]string, []error) {
		return validation.StringInSlice(valid, ignoreCase)(i, k)
	}
}
