// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package validators

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ validator.String = MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator{}
)

type MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator struct {
	OtherFieldExpression path.Expression
	OtherFieldValueRegex *regexp.Regexp
	ErrorMessage         string
}

type MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (av MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator) MarkdownDescription(_ context.Context) string {
	return av.ErrorMessage
}

func (av MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidator) Validate(ctx context.Context, req MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorRequest, res *MakeFieldRequiredWhenOtherFieldDoesNotHaveValueValidatorResponse) {
	paths, _ := req.Config.PathMatches(ctx, av.OtherFieldExpression)
	if paths == nil || len(paths) != 1 {
		res.Diagnostics.AddError(
			"Validator Configuration Error: Other field match failed",
			"The validator could not uniquely locate the other field in the configuration. Ensure that 'OtherFieldExpression' matches exactly one attribute.",
		)
		return
	}
	otherFieldValue := ""
	diags := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
	isUnknown := false
	if diags.HasError() {
		isUnknown = strings.Contains(diags.Errors()[0].Detail(), "Received unknown value") || strings.Contains(diags.Errors()[0].Summary(), "Received unknown value")
	}

	if req.ConfigValue.IsUnknown() || isUnknown {
		tflog.Debug(ctx, "RequiredFieldValidator: Value is unknown, skipping validation")
		return
	}

	doesNotMatchCorrectly := !av.OtherFieldValueRegex.MatchString(otherFieldValue)
	currentValueNotDefined := req.ConfigValue.IsNull()

	if (doesNotMatchCorrectly && !currentValueNotDefined) || (!doesNotMatchCorrectly && currentValueNotDefined) {
		res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
		return
	}
}
