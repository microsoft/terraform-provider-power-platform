// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DynamicColumns returns a ConfigValidator that ensures that the given expression
// many-to-one relationships are using set collections.

// OtherFieldRequiredWhenValueOfValidator is a validator that ensures that a field is required when the value of another field is equal to a specific value.
func OtherFieldRequiredWhenValueOf(otherFieldExpression path.Expression, currentFieldValueRegex *regexp.Regexp, errorMessage string) validator.String {
	return &OtherFieldRequiredWhenValueOfValidator{
		OtherFieldExpression:   otherFieldExpression,
		CurrentFieldValueRegex: currentFieldValueRegex,
		ErrorMessage:           errorMessage,
	}
}
