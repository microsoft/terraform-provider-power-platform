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
)

var (
	_ validator.String = OtherFieldRequiredWhenValueOfValidator{}
)

type OtherFieldRequiredWhenValueOfValidator struct {
	OtherFieldExpression   path.Expression
	OtherFieldValueRegex   *regexp.Regexp
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
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	currentFieldValue := ""
	diags := req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	if (av.CurrentFieldValueRegex != nil && av.CurrentFieldValueRegex.MatchString(currentFieldValue)) || (av.CurrentFieldValueRegex == nil && currentFieldValue != "") {
		paths, diags := req.Config.PathMatches(ctx, av.OtherFieldExpression)
		res.Diagnostics.Append(diags...)
		if res.Diagnostics.HasError() {
			return
		}
		if paths == nil || len(paths) != 1 {
			res.Diagnostics.AddError("Other field required when value of validator should have exactly one match", "")
			return
		}

		otherFieldValue := new(string)
		d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)

		isUnknown := false
		if d.HasError() {
			isUnknown = strings.Contains(d.Errors()[0].Detail(), "Received unknown value") || strings.Contains(d.Errors()[0].Summary(), "Received unknown value")
		}

		if (av.OtherFieldValueRegex != nil && otherFieldValue != nil && !av.OtherFieldValueRegex.MatchString(*otherFieldValue)) ||
			(av.OtherFieldValueRegex == nil && (otherFieldValue == nil || *otherFieldValue == "") && !isUnknown) {
			res.Diagnostics.AddError(
				av.ErrorMessage,
				"Field \""+paths[0].String()+"\" does not meet required value conditions.",
			)
		}
	}
}
