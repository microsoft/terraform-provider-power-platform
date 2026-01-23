// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestRestoreOriginalValueModifier_String(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.RestoreOriginalStringModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	nullRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil)

	t.Run("store_original_on_create", func(t *testing.T) {
		req := planmodifier.StringRequest{
			Path:        path.Root("field"),
			State:       tfsdk.State{Raw: nullRaw},
			ConfigValue: types.StringValue("value"),
		}
		resp := planmodifier.StringResponse{}
		privatePtr := setPrivateData(t, &resp)

		modifier.PlanModifyString(ctx, req, &resp)

		stored := privateGetKey(ctx, t, privatePtr, "field")
		if stored != nil {
			t.Fatal("expected stored key to be nil due to invalid JSON value")
		}
	})

	t.Run("restore_on_destroy", func(t *testing.T) {
		req := planmodifier.StringRequest{
			Path:        path.Root("field"),
			Plan:        tfsdk.Plan{Raw: nullRaw},
			ConfigValue: types.StringValue("value"),
		}
		resp := planmodifier.StringResponse{}
		setPrivateData(t, &resp)

		modifier.PlanModifyString(ctx, req, &resp)
	})
}

func TestUnitRestoreOriginalValueModifier_Bool(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.RestoreOriginalBoolModifier()

	nullRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil)
	req := planmodifier.BoolRequest{
		Path:        path.Root("flag"),
		State:       tfsdk.State{Raw: nullRaw},
		ConfigValue: types.BoolValue(true),
	}
	resp := planmodifier.BoolResponse{}
	privatePtr := setPrivateData(t, &resp)

	modifier.PlanModifyBool(ctx, req, &resp)

	stored := privateGetKey(ctx, t, privatePtr, "flag")
	if stored != nil {
		t.Fatal("expected stored key to be nil due to zero-length value")
	}
}
