// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// OtherFieldRequiredWhenValueOfValidator is a validator that ensures that a field is required (with optional specific value `currentFieldValueRegex`) when the value of another field is equal to a value.
// When otherFieldValueRegex matches the value of the field specified by OtherFieldExpression, the other field will have to have the required value specified.
func OtherFieldRequiredWhenValueOf(otherFieldExpression path.Expression, otherFieldValueRegex, currentFieldValueRegex *regexp.Regexp, errorMessage string) validator.String {
	return &OtherFieldRequiredWhenValueOfValidator{
		OtherFieldExpression:   otherFieldExpression,
		OtherFieldValueRegex:   otherFieldValueRegex,
		CurrentFieldValueRegex: currentFieldValueRegex,
		ErrorMessage:           errorMessage,
	}
}
