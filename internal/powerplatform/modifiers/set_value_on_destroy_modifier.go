// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_modifiers

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SetValueOnDestroyStringModifier(v string) planmodifier.String {
	return &setValueOnDestroyModifier{
		value: []byte(v),
	}
}

func SetValueOnDestroyBoolModifier(v bool) planmodifier.Bool {
	return &setValueOnDestroyModifier{
		value: func() []byte {
			if v {
				return []byte{1}
			}
			return []byte{0}
		}(),
	}
}

type setValueOnDestroyModifier struct {
	value []byte
}

func (d *setValueOnDestroyModifier) Description(ctx context.Context) string {
	return "Stores the original value of an attribute that can't be destroyed so that it can be set to its original value when the resource is destroyed."
}

func (d *setValueOnDestroyModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *setValueOnDestroyModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Check if the resource is being destroyed
	if req.Plan.Raw.IsNull() {
		if !req.ConfigValue.IsNull() {
			log.Default().Printf("Restoring original value for attribute %s", req.PathExpression.String())
			resp.PlanValue = basetypes.NewStringValue(string(d.value))
		}
		//req.Plan.SetAttribute(ctx, req.Path, req.Private.GetKey(ctx, req.Path.String()))
	}
}

func (d *setValueOnDestroyModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
}
