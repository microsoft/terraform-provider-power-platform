// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestSetStringAttributeUnknownOnlyIfSecondAttributeChange(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.SetStringAttributeUnknownOnlyIfSecondAttributeChange(path.Root("second"))

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	schemaDef := schema.Schema{Attributes: map[string]schema.Attribute{
		"second": schema.StringAttribute{Optional: true},
	}}

	req := planmodifier.StringRequest{
		Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
			"second": tftypes.NewValue(tftypes.String, "new"),
		}),
		State: newState(t, schemaDef, map[string]tftypes.Value{
			"second": tftypes.NewValue(tftypes.String, "old"),
		}),
		PlanValue: types.StringValue("keep"),
	}
	resp := planmodifier.StringResponse{}

	modifier.PlanModifyString(ctx, req, &resp)

	if !resp.PlanValue.IsUnknown() {
		t.Fatal("expected plan value to be unknown")
	}
}
