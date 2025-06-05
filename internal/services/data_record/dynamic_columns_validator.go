// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ resource.ConfigValidator = &DynamicsColumnsValidator{}

type DynamicsColumnsValidator struct {
	PathExpression path.Expression
}

func (v DynamicsColumnsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v DynamicsColumnsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Dynamic columns should use set for many-to-one relationships: %s", v.PathExpression)
}

func (v DynamicsColumnsValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v DynamicsColumnsValidator) ValidateProvider(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v DynamicsColumnsValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v DynamicsColumnsValidator) ValidateEphemeralResource(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v DynamicsColumnsValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var diags diag.Diagnostics

	matchedPaths, matchedPathsDiags := config.PathMatches(ctx, v.PathExpression)
	diags.Append(matchedPathsDiags...)

	if matchedPaths == nil || len(matchedPaths) != 1 {
		diags.AddError("Dynamic columns validator shoud have exactly one match", "")
		return diags
	}

	var dynamicColumns types.Dynamic
	if err := config.GetAttribute(ctx, matchedPaths[0], &dynamicColumns); err != nil {
		diags.AddError("Failed to get dynamic columns attribute", err.Error())
		return diags
	}

	terraformDynamicColumns, err := dynamicColumns.UnderlyingValue().ToTerraformValue(ctx)
	if err != nil {
		diags.AddError("Failed to convert dynamic columns to terraform value", err.Error())
		return diags
	}

	var attrs map[string]tftypes.Value
	err = terraformDynamicColumns.As(&attrs)
	if err != nil {
		diags.AddError("Failed to convert dynamic columns to map[string]tftypes.Value", err.Error())
		return diags
	}

	for key, value := range attrs {
		if value.Type().Is(tftypes.Tuple{}) || value.Type().Is(tftypes.List{}) {
			msg := fmt.Sprintf("Dynamic columns should use set collection with `toset([...])` for many-to-one relationships. Record attribute: '%s'", key)
			diags.AddWarning(msg, msg)
		}
	}

	return diags
}
