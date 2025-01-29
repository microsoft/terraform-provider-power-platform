// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	_ validator.String = OtherFieldRequiredWhenValueOfValidator{}
)

type OtherFieldRequiredWhenValueOfValidator struct {
	OtherFieldExpression   path.Expression
	CurrentFieldValueRegex *regexp.Regexp
	ErrorMessage           string
}

type OtherFieldRequiredWhenValueOfValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type OtherFieldRequiredWhenValueOfValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (av OtherFieldRequiredWhenValueOfValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av OtherFieldRequiredWhenValueOfValidator) MarkdownDescription(_ context.Context) string {
	return av.ErrorMessage
}

func (av OtherFieldRequiredWhenValueOfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := OtherFieldRequiredWhenValueOfValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &OtherFieldRequiredWhenValueOfValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av OtherFieldRequiredWhenValueOfValidator) Validate(ctx context.Context, req OtherFieldRequiredWhenValueOfValidatorRequest, res *OtherFieldRequiredWhenValueOfValidatorResponse) {
	// If attribute configuration is null, there is nothing else to validate
	// If attribute configuration is unknown, delay the validation until it is known.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	currentFieldValue := ""
	_ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)

	if av.CurrentFieldValueRegex.MatchString(currentFieldValue) {
		paths, _ := req.Config.PathMatches(ctx, av.OtherFieldExpression)
		if paths == nil && len(paths) != 1 {
			res.Diagnostics.AddError("Other field required when value of validator should have exactly one match", "")
			return
		}

		otherFieldValue := ""
		_ = req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
		if otherFieldValue == "" {
			res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
		}
	}
}
