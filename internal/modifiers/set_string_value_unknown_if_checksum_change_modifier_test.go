// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func TestUnitSetStringValueToUnknownIfChecksumsChangeModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.SetStringValueToUnknownIfChecksumsChangeModifier(
		[]string{"file1", "checksum1"},
		[]string{"file2", "checksum2"},
	)

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	t.Run("checksum_change_sets_unknown", func(t *testing.T) {
		tmpDir := t.TempDir()
		file1 := filepath.Join(tmpDir, "file1.txt")
		file2 := filepath.Join(tmpDir, "file2.txt")
		if err := os.WriteFile(file1, []byte("one"), 0600); err != nil {
			t.Fatalf("write file1: %v", err)
		}
		if err := os.WriteFile(file2, []byte("two"), 0600); err != nil {
			t.Fatalf("write file2: %v", err)
		}

		checksum2, err := helpers.CalculateSHA256(file2)
		if err != nil {
			t.Fatalf("checksum2: %v", err)
		}

		planSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"file1": schema.StringAttribute{Optional: true},
			"file2": schema.StringAttribute{Optional: true},
		}}
		stateSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"checksum1": schema.StringAttribute{Optional: true},
			"checksum2": schema.StringAttribute{Optional: true},
		}}

		req := planmodifier.StringRequest{
			Plan: newPlan(t, planSchema, map[string]tftypes.Value{
				"file1": tftypes.NewValue(tftypes.String, file1),
				"file2": tftypes.NewValue(tftypes.String, file2),
			}),
			State: newState(t, stateSchema, map[string]tftypes.Value{
				"checksum1": tftypes.NewValue(tftypes.String, "different"),
				"checksum2": tftypes.NewValue(tftypes.String, checksum2),
			}),
			PlanValue: types.StringValue("keep"),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.PlanValue.IsUnknown() {
			t.Fatal("expected plan value to be unknown")
		}
	})

	t.Run("plan_attribute_missing", func(t *testing.T) {
		planSchema := schema.Schema{Attributes: map[string]schema.Attribute{}}
		stateSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"checksum1": schema.StringAttribute{Optional: true},
		}}

		req := planmodifier.StringRequest{
			Plan: newPlan(t, planSchema, map[string]tftypes.Value{}),
			State: newState(t, stateSchema, map[string]tftypes.Value{
				"checksum1": tftypes.NewValue(tftypes.String, ""),
			}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics to have error")
		}
	})

	t.Run("state_attribute_missing", func(t *testing.T) {
		planSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"file1": schema.StringAttribute{Optional: true},
		}}
		stateSchema := schema.Schema{Attributes: map[string]schema.Attribute{}}

		req := planmodifier.StringRequest{
			Plan: newPlan(t, planSchema, map[string]tftypes.Value{
				"file1": tftypes.NewValue(tftypes.String, ""),
			}),
			State: newState(t, stateSchema, map[string]tftypes.Value{}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics to have error")
		}
	})

	t.Run("checksum_error_adds_diagnostic", func(t *testing.T) {
		tmpDir := t.TempDir()
		planSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"file1": schema.StringAttribute{Optional: true},
			"file2": schema.StringAttribute{Optional: true},
		}}
		stateSchema := schema.Schema{Attributes: map[string]schema.Attribute{
			"checksum1": schema.StringAttribute{Optional: true},
			"checksum2": schema.StringAttribute{Optional: true},
		}}

		req := planmodifier.StringRequest{
			Plan: newPlan(t, planSchema, map[string]tftypes.Value{
				"file1": tftypes.NewValue(tftypes.String, tmpDir),
				"file2": tftypes.NewValue(tftypes.String, ""),
			}),
			State: newState(t, stateSchema, map[string]tftypes.Value{
				"checksum1": tftypes.NewValue(tftypes.String, ""),
				"checksum2": tftypes.NewValue(tftypes.String, ""),
			}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics to have error")
		}
	})
}
