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

func TestUnitUseStateForUnknownKeepNonNullStateModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.UseStateForUnknownKeepNonNullStateModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("state_null", func(t *testing.T) {
		req := planmodifier.StringRequest{
			StateValue:  types.StringNull(),
			PlanValue:   types.StringUnknown(),
			ConfigValue: types.StringValue("config"),
		}
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.PlanValue.ValueString() != "keep" {
			t.Fatal("expected plan value to remain unchanged")
		}
	})

	t.Run("plan_known", func(t *testing.T) {
		req := planmodifier.StringRequest{
			StateValue:  types.StringValue("state"),
			PlanValue:   types.StringValue("plan"),
			ConfigValue: types.StringValue("config"),
		}
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.PlanValue.ValueString() != "keep" {
			t.Fatal("expected plan value to remain unchanged")
		}
	})

	t.Run("config_unknown", func(t *testing.T) {
		req := planmodifier.StringRequest{
			StateValue:  types.StringValue("state"),
			PlanValue:   types.StringUnknown(),
			ConfigValue: types.StringUnknown(),
		}
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.PlanValue.ValueString() != "keep" {
			t.Fatal("expected plan value to remain unchanged")
		}
	})

	t.Run("use_state_for_unknown", func(t *testing.T) {
		req := planmodifier.StringRequest{
			StateValue:  types.StringValue("state"),
			PlanValue:   types.StringUnknown(),
			ConfigValue: types.StringValue("config"),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.PlanValue.ValueString() != "state" {
			t.Fatal("expected plan value to be state value")
		}
	})
}
