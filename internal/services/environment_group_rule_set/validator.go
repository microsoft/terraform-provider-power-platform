// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

// this code extends https://github.com/hashicorp/terraform-plugin-framework-validators/blob/8beb21856fdd36e5f84825e78fe10fa49212c41d/internal/schemavalidator/also_requires.go

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

// This type of validator must satisfy all types.
var (
	_ validator.Bool    = AlsoRequiresValueValidator{}
	_ validator.Float32 = AlsoRequiresValueValidator{}
	_ validator.Float64 = AlsoRequiresValueValidator{}
	_ validator.Int32   = AlsoRequiresValueValidator{}
	_ validator.Int64   = AlsoRequiresValueValidator{}
	_ validator.List    = AlsoRequiresValueValidator{}
	_ validator.Map     = AlsoRequiresValueValidator{}
	_ validator.Number  = AlsoRequiresValueValidator{}
	_ validator.Object  = AlsoRequiresValueValidator{}
	_ validator.Set     = AlsoRequiresValueValidator{}
	_ validator.String  = AlsoRequiresValueValidator{}
)

// AlsoRequiresValueValidator is the underlying struct implementing AlsoRequires.
type AlsoRequiresValueValidator struct {
	PathExpressions path.Expressions
	ExpectedValue   attr.Value
}

type AlsoRequiresValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type AlsoRequiresValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (av AlsoRequiresValueValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av AlsoRequiresValueValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set with a given value, also these are set: %q", av.PathExpressions)
}

func (av AlsoRequiresValueValidator) Validate(ctx context.Context, req AlsoRequiresValueValidatorRequest, res *AlsoRequiresValueValidatorResponse) {
	// If attribute configuration is null, there is nothing else to validate
	if req.ConfigValue.IsNull() {
		return
	}

	expressions := req.PathExpression.MergeExpressions(av.PathExpressions...)

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		res.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.Path) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if mpVal.IsUnknown() {
				return
			}

			if mpVal.IsNull() {
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q must be specified when %q is specified", mp, req.Path),
				))
			}

			if !mpVal.Equal(av.ExpectedValue) {
				executedFromAttrName, _ := req.Path.Steps().LastStep()
				executedForAttrName, _ := mp.Steps().LastStep()
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q did not had the expected value of %q but %q. This is requried by %q", executedForAttrName.String(), av.ExpectedValue.String(), mpVal.String(), executedFromAttrName.String()),
				))
			}
		}
	}
}

func (av AlsoRequiresValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av AlsoRequiresValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := AlsoRequiresValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &AlsoRequiresValueValidatorResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func AlsoRequiresValueString(expectedValue attr.Value, expressions ...path.Expression) validator.String {
	return AlsoRequiresValueValidator{
		PathExpressions: expressions,
		ExpectedValue:   expectedValue,
	}
}

func AlsoRequiresValueInt32(expectedValue attr.Value, expressions ...path.Expression) validator.Int32 {
	return AlsoRequiresValueValidator{
		PathExpressions: expressions,
		ExpectedValue:   expectedValue,
	}
}

func AlsoRequiresValueBool(expectedValue attr.Value, expressions ...path.Expression) validator.Bool {
	return AlsoRequiresValueValidator{
		PathExpressions: expressions,
		ExpectedValue:   expectedValue,
	}
}
