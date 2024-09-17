// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequireReplaceObjectToEmptyModifier() planmodifier.Object {
	return &requireReplaceObjectToEmptyModifier{}
}

type requireReplaceObjectToEmptyModifier struct {
}

func (d *requireReplaceObjectToEmptyModifier) Description(ctx context.Context) string {
	return "Ensures that change to empty attribute value will force a replace when changed."
}

func (d *requireReplaceObjectToEmptyModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *requireReplaceObjectToEmptyModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.StateValue.IsNull() && req.PlanValue.IsNull() {
		return
	}

	//we only replace is object was created and is being set to empty/nil now
	if !req.StateValue.IsNull() && (req.PlanValue.Attributes() == nil || len(req.PlanValue.Attributes()) == 0) {
		resp.RequiresReplace = true
	}
}
