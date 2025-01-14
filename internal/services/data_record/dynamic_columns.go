// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// DynamicColumns returns a ConfigValidator that ensures that the given expression
// many-to-one relationships are using set collections
func DynamicColumns(expression path.Expression) resource.ConfigValidator {
	return &DynamicsColumnsValidator{
		PathExpression: expression,
	}
}
