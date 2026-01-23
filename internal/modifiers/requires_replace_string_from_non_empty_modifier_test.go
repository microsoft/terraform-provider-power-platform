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

func TestUnitRequireReplaceStringFromNonEmptyPlanModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.RequireReplaceStringFromNonEmptyPlanModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("requires_replace", func(t *testing.T) {
		req := planmodifier.StringRequest{
			PlanValue:  types.StringValue("new"),
			StateValue: types.StringValue("old"),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be true")
		}
	})

	t.Run("no_replace_for_empty", func(t *testing.T) {
		req := planmodifier.StringRequest{
			PlanValue:  types.StringValue("new"),
			StateValue: types.StringValue(""),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.RequiresReplace {
			t.Fatal("expected RequiresReplace to be false")
		}
	})
}
