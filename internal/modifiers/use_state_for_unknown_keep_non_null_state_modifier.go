// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

// terraform-plugin-framework v1.15.1 changed behaviour of UseStateForUnknown plan modifier.
// Before the change, if new nested map, that contains Computed/Unknown elements is added to already created resource,
// those elements are treated as Unknown.
// After the change, those elements are treated as Null.
// This breaks bulk operations, as the ID of newly created objects can never be set, as Terraform would expects them to be Null.
// Ref: https://github.com/hashicorp/terraform-plugin-framework/issues/1211

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForUnknownKeepNonNullStateModifier() planmodifier.String {
	return &useStateForUnknownKeepNonNullStateModifier{}
}

type useStateForUnknownKeepNonNullStateModifier struct {
}

func (d *useStateForUnknownKeepNonNullStateModifier) Description(ctx context.Context) string {
	return "Once set to a non-null value, the value of this attribute in state will not change."
}

func (d *useStateForUnknownKeepNonNullStateModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *useStateForUnknownKeepNonNullStateModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there the state value is null.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
