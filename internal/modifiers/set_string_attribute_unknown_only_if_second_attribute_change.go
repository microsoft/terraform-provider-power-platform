// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetStringAttributeUnknownOnlyIfSecondAttributeChange(secondAttributePath path.Path) planmodifier.String {
	return &setStringAttributeUnknownOnlyIfSecondAttributeChange{
		secondAttributePath: secondAttributePath,
	}
}

type setStringAttributeUnknownOnlyIfSecondAttributeChange struct {
	secondAttributePath path.Path
}

func (d *setStringAttributeUnknownOnlyIfSecondAttributeChange) Description(ctx context.Context) string {
	return "Ensures that attribute is set to unknown, only if second attribute changes."
}

func (d *setStringAttributeUnknownOnlyIfSecondAttributeChange) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *setStringAttributeUnknownOnlyIfSecondAttributeChange) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planSecondAttribute types.String
	diags := req.Plan.GetAttribute(ctx, d.secondAttributePath, &planSecondAttribute)
	resp.Diagnostics.Append(diags...)

	var stateSecondAttribute types.String
	diags = req.State.GetAttribute(ctx, d.secondAttributePath, &stateSecondAttribute)
	resp.Diagnostics.Append(diags...)

	if planSecondAttribute.ValueString() != stateSecondAttribute.ValueString() && !planSecondAttribute.IsUnknown() && !planSecondAttribute.IsNull() {
		resp.PlanValue = types.StringUnknown()
	}
}
