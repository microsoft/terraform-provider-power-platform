// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseStateForUnknownWhenParentInStateListModifier copies the prior state value (including null) of a nested
// list attribute into an unknown plan value, but only when the parent object already exists in state.
// When the parent object is being added to an existing resource the value is left unknown, so the
// provider can still set whatever the service returns after apply.
func UseStateForUnknownWhenParentInStateListModifier() planmodifier.List {
	return &useStateForUnknownWhenParentInStateListModifier{}
}

type useStateForUnknownWhenParentInStateListModifier struct {
}

func (d *useStateForUnknownWhenParentInStateListModifier) Description(ctx context.Context) string {
	return "Once the parent object exists in state, the value of this attribute in state will not change."
}

func (d *useStateForUnknownWhenParentInStateListModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *useStateForUnknownWhenParentInStateListModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Do nothing if there is no state (resource is being created).
	if req.State.Raw.IsNull() {
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

	var parent types.Object
	diags := req.State.GetAttribute(ctx, req.Path.ParentPath(), &parent)
	if diags.HasError() || parent.IsNull() || parent.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
