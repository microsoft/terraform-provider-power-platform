// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseDefaultBoolValueWhenConfigNull(defaultValue bool) planmodifier.Bool {
	return &useDefaultBoolValueWhenConfigNullModifier{
		defaultValue: defaultValue,
	}
}

type useDefaultBoolValueWhenConfigNullModifier struct {
	defaultValue bool
}

func (d *useDefaultBoolValueWhenConfigNullModifier) Description(ctx context.Context) string {
	return "When configuration omits this attribute, use the resource default value instead of inheriting the current state value."
}

func (d *useDefaultBoolValueWhenConfigNullModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *useDefaultBoolValueWhenConfigNullModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsUnknown() || !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = types.BoolValue(d.defaultValue)
}
