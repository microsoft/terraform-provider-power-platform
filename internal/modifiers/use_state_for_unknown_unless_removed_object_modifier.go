// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknownUnlessRemovedObjectModifier copies the prior state value (including null) of an
// Optional+Computed object into an unknown plan value, except when the object is being removed from the
// configuration. In that case the plan value is left unknown, so a RequiresReplace raised for the removal is
// still honored: Terraform ignores RequiresReplace when the planned value equals the prior state value.
func UseStateForUnknownUnlessRemovedObjectModifier() planmodifier.Object {
	return &useStateForUnknownUnlessRemovedObjectModifier{}
}

type useStateForUnknownUnlessRemovedObjectModifier struct {
}

func (d *useStateForUnknownUnlessRemovedObjectModifier) Description(ctx context.Context) string {
	return "Unless the object is removed from the configuration, the value of this attribute in state will not change."
}

func (d *useStateForUnknownUnlessRemovedObjectModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *useStateForUnknownUnlessRemovedObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
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

	// Do nothing if the object is being removed from the configuration.
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}
