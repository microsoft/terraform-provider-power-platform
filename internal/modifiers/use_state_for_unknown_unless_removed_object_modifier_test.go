// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestUnitUseStateForUnknownUnlessRemovedObjectModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.UseStateForUnknownUnlessRemovedObjectModifier()

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	sch := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Optional: true},
		},
	}
	existingResource := newState(t, sch, map[string]tftypes.Value{"name": tftypes.NewValue(tftypes.String, "name")})
	newResource := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)}

	attrTypes := map[string]attr.Type{"value": types.StringType}
	objectValue := func(value string) types.Object {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{"value": types.StringValue(value)})
	}
	unknown := types.ObjectUnknown(attrTypes)
	null := types.ObjectNull(attrTypes)

	run := func(state tfsdk.State, stateValue, planValue, configValue types.Object) types.Object {
		req := planmodifier.ObjectRequest{
			State:       state,
			StateValue:  stateValue,
			PlanValue:   planValue,
			ConfigValue: configValue,
		}
		resp := planmodifier.ObjectResponse{PlanValue: planValue}
		modifier.PlanModifyObject(ctx, req, &resp)
		return resp.PlanValue
	}

	t.Run("resource_created", func(t *testing.T) {
		if got := run(newResource, null, unknown, null); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown")
		}
	})

	t.Run("plan_known", func(t *testing.T) {
		if got := run(existingResource, objectValue("state"), objectValue("plan"), objectValue("plan")); !got.Equal(objectValue("plan")) {
			t.Fatal("expected plan value to remain unchanged")
		}
	})

	t.Run("config_unknown", func(t *testing.T) {
		if got := run(existingResource, objectValue("state"), unknown, unknown); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown")
		}
	})

	t.Run("object_removed_from_config", func(t *testing.T) {
		if got := run(existingResource, objectValue("state"), unknown, null); !got.IsUnknown() {
			t.Fatal("expected plan value to remain unknown when the object is removed from the configuration")
		}
	})

	t.Run("use_state_for_unknown", func(t *testing.T) {
		if got := run(existingResource, objectValue("state"), unknown, objectValue("config")); !got.Equal(objectValue("state")) {
			t.Fatalf("expected plan value to use state value, got %s", got)
		}
	})

	t.Run("use_null_state_for_unknown", func(t *testing.T) {
		if got := run(existingResource, null, unknown, null); !got.IsNull() {
			t.Fatalf("expected plan value to use null state value, got %s", got)
		}
	})
}
