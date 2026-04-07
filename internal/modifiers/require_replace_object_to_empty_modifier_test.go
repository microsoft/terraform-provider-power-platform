// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestUnitRequireReplaceObjectToEmptyModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.RequireReplaceObjectToEmptyModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("both_null", func(t *testing.T) {
		attrTypes := map[string]attr.Type{}
		req := planmodifier.ObjectRequest{
			PlanValue:  types.ObjectNull(attrTypes),
			StateValue: types.ObjectNull(attrTypes),
		}
		resp := planmodifier.ObjectResponse{}

		modifier.PlanModifyObject(ctx, req, &resp)

		if resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be false")
		}
	})

	t.Run("state_non_null_plan_empty", func(t *testing.T) {
		attrTypes := map[string]attr.Type{}
		stateValue, diags := types.ObjectValue(attrTypes, map[string]attr.Value{})
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		planValue, diags := types.ObjectValue(attrTypes, map[string]attr.Value{})
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		req := planmodifier.ObjectRequest{
			PlanValue:  planValue,
			StateValue: stateValue,
		}
		resp := planmodifier.ObjectResponse{}

		modifier.PlanModifyObject(ctx, req, &resp)

		if !resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be true")
		}
	})
}
