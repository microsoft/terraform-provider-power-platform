// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RestoreOriginalStringModifier() planmodifier.String {
	return &restoreOriginalValueModifier{}
}

func RestoreOriginalBoolModifier() planmodifier.Bool {
	return &restoreOriginalValueModifier{}
}

type restoreOriginalValueModifier struct {
}

func (d *restoreOriginalValueModifier) Description(ctx context.Context) string {
	return "Stores the original value of an attribute that can't be destroyed so that it can be set to its original value when the resource is destroyed."
}

func (d *restoreOriginalValueModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *restoreOriginalValueModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Check if the resource is being created
	if req.State.Raw.IsNull() {
		if !req.ConfigValue.IsNull() {
			log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
			resp.Private.SetKey(ctx, req.Path.String(), []byte{1})
		}
	}

	// Check if the resource is being destroyed
	if req.Plan.Raw.IsNull() {
		if !req.ConfigValue.IsNull() {
			log.Default().Printf("Restoring original value for attribute %s", req.PathExpression.String())
		}
	}
}

func (d *restoreOriginalValueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Check if the resource is being created
	if req.State.Raw.IsNull() {
		if !req.ConfigValue.IsNull() {
			log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
			resp.Private.SetKey(ctx, req.Path.String(), []byte{})
		}
	}
}
