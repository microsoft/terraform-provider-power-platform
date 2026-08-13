// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestUnitUseDefaultBoolValueWhenConfigNullModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.UseDefaultBoolValueWhenConfigNull(false)

	t.Run("uses default when config is null", func(t *testing.T) {
		req := planmodifier.BoolRequest{
			ConfigValue: types.BoolNull(),
			StateValue:  types.BoolValue(true),
			PlanValue:   types.BoolValue(true),
		}
		resp := planmodifier.BoolResponse{}

		modifier.PlanModifyBool(ctx, req, &resp)

		if resp.PlanValue.IsNull() || resp.PlanValue.IsUnknown() {
			t.Fatal("expected known plan value")
		}

		if resp.PlanValue.ValueBool() {
			t.Fatal("expected plan value false, got true")
		}
	})

	t.Run("does not override explicit config", func(t *testing.T) {
		req := planmodifier.BoolRequest{
			ConfigValue: types.BoolValue(true),
			StateValue:  types.BoolValue(false),
			PlanValue:   types.BoolValue(true),
		}
		resp := planmodifier.BoolResponse{PlanValue: types.BoolValue(true)}

		modifier.PlanModifyBool(ctx, req, &resp)

		if !resp.PlanValue.ValueBool() {
			t.Fatal("expected explicit configured value to be preserved")
		}
	})
}
