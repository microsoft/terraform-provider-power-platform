// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MakeFieldRequiredWhenOtherFieldDoesNotHaveValue(otherFieldExpression path.Expression, otherFieldValueRegex *regexp.Regexp, errorMessage string) validator.String {
	return &MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator{
		OtherFieldExpression: otherFieldExpression,
		OtherFieldValueRegex: otherFieldValueRegex,
		ErrorMessage:         errorMessage,
	}
}
