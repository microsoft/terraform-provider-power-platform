// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestUnitUseStateForUnknownWhenParentInStateListModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.UseStateForUnknownWhenParentInStateListModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	sch := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"parent": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"items": schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType},
				},
			},
		},
	}
	listType := tftypes.List{ElementType: tftypes.String}
	parentType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"items": listType}}
	itemsPath := path.Root("parent").AtName("items")

	stateWithParent := newState(t, sch, map[string]tftypes.Value{
		"parent": tftypes.NewValue(parentType, map[string]tftypes.Value{"items": tftypes.NewValue(listType, nil)}),
	})
	stateWithoutParent := newState(t, sch, map[string]tftypes.Value{
		"parent": tftypes.NewValue(parentType, nil),
	})
	stateValue := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")})

	run := func(state tfsdk.State, stateValue, planValue, configValue types.List) types.List {
		req := planmodifier.ListRequest{
			Path:        itemsPath,
			State:       state,
			StateValue:  stateValue,
			PlanValue:   planValue,
			ConfigValue: configValue,
		}
		resp := planmodifier.ListResponse{PlanValue: planValue}
		modifier.PlanModifyList(ctx, req, &resp)
		return resp.PlanValue
	}

	t.Run("resource_created", func(t *testing.T) {
		state := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)}
		if got := run(state, types.ListNull(types.StringType), types.ListUnknown(types.StringType), types.ListNull(types.StringType)); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown")
		}
	})

	t.Run("plan_known", func(t *testing.T) {
		planValue := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("plan")})
		if got := run(stateWithParent, stateValue, planValue, planValue); !got.Equal(planValue) {
			t.Fatal("expected plan value to remain unchanged")
		}
	})

	t.Run("config_unknown", func(t *testing.T) {
		if got := run(stateWithParent, stateValue, types.ListUnknown(types.StringType), types.ListUnknown(types.StringType)); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown")
		}
	})

	t.Run("parent_not_in_state", func(t *testing.T) {
		if got := run(stateWithoutParent, types.ListNull(types.StringType), types.ListUnknown(types.StringType), types.ListNull(types.StringType)); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown when the parent is being added")
		}
	})

	t.Run("use_state_for_unknown", func(t *testing.T) {
		if got := run(stateWithParent, stateValue, types.ListUnknown(types.StringType), types.ListNull(types.StringType)); !got.Equal(stateValue) {
			t.Fatalf("expected plan value to use state value, got %s", got)
		}
	})

	t.Run("use_null_state_for_unknown", func(t *testing.T) {
		if got := run(stateWithParent, types.ListNull(types.StringType), types.ListUnknown(types.StringType), types.ListNull(types.StringType)); !got.IsNull() {
			t.Fatalf("expected plan value to use null state value, got %s", got)
		}
	})
}
