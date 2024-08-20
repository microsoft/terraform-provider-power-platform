// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequireReplaceStringFromNonEmptyPlanModifier() planmodifier.String {
	return &requireReplaceStringFromNonEmptyPlanModifier{}
}

type requireReplaceStringFromNonEmptyPlanModifier struct {
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) Description(ctx context.Context) string {
	return "Ensures that change from non empty attribute value will force a replace when changed."
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueString() != "") {
		resp.RequiresReplace = true
	}
}
