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

func TestRequireReplaceIntAttributePlanModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.RequireReplaceIntAttributePlanModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("requires_replace", func(t *testing.T) {
		req := planmodifier.Int64Request{
			PlanValue:  types.Int64Value(10),
			StateValue: types.Int64Value(5),
		}
		resp := planmodifier.Int64Response{}

		modifier.PlanModifyInt64(ctx, req, &resp)

		if !resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be true")
		}
	})

	t.Run("no_replace_for_zero", func(t *testing.T) {
		req := planmodifier.Int64Request{
			PlanValue:  types.Int64Value(10),
			StateValue: types.Int64Value(0),
		}
		resp := planmodifier.Int64Response{}

		modifier.PlanModifyInt64(ctx, req, &resp)

		if resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be false")
		}
	})
}
