// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func SyncAttributePlanModifier(syncAttribute string) planmodifier.String {
	return &syncAttributePlanModifier{
		syncAttribute: syncAttribute,
	}
}

type syncAttributePlanModifier struct {
	syncAttribute string
}

func (d *syncAttributePlanModifier) Description(ctx context.Context) string {
	return "Ensures that file attribute and file checksum attribute are kept synchronised."
}

func (d *syncAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *syncAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var settingsFile types.String
	diags := req.Plan.GetAttribute(ctx, path.Root(d.syncAttribute), &settingsFile)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if settingsFile.IsNull() {
		resp.PlanValue = types.StringNull()
	} else if settingsFile.IsUnknown() {
		resp.PlanValue = types.StringNull()
	} else {
		value, err := helpers.CalculateSHA256(settingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())
			return
		}

		if value == "" {
			resp.PlanValue = types.StringUnknown()
		} else {
			resp.PlanValue = types.StringValue(value)
		}
	}
}
