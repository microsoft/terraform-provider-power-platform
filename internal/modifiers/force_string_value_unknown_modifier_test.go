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

func TestForceStringValueUnknownModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.ForceStringValueUnknownModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("no_key", func(t *testing.T) {
		req := planmodifier.StringRequest{}
		setPrivateData(t, &req)
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.PlanValue.IsNull() {
			t.Fatal("expected plan value to remain null")
		}
	})

	t.Run("force_unknown", func(t *testing.T) {
		req := planmodifier.StringRequest{}
		privatePtr := setPrivateData(t, &req)
		privateSetKey(ctx, t, privatePtr, "force_value_unknown", []byte("true"))
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.PlanValue.IsUnknown() {
			t.Fatal("expected plan value to be unknown")
		}
	})
}
