// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequireReplaceIntAttributePlanModifier() planmodifier.Int64 {
	return &requireReplaceIntAttributePlanModifier{}
}

type requireReplaceIntAttributePlanModifier struct {
}

func (d *requireReplaceIntAttributePlanModifier) Description(ctx context.Context) string {
	return "Ensures that int attribute will force a replace when changed."
}

func (d *requireReplaceIntAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *requireReplaceIntAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.PlanValue != req.StateValue {
		resp.RequiresReplace = true
	}
}