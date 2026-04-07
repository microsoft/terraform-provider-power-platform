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

func TestUnitSyncAttributePlanModifier(t *testing.T) {
	ctx := context.Background()
	modifier := modifiers.SyncAttributePlanModifier("file")

	if desc := modifier.(interface{ Description(context.Context) string }).Description(ctx); desc == "" {
		t.Fatal("expected Description to be non-empty")
	}
	if desc := modifier.(interface{ MarkdownDescription(context.Context) string }).MarkdownDescription(ctx); desc == "" {
		t.Fatal("expected MarkdownDescription to be non-empty")
	}

	schemaDef := schema.Schema{Attributes: map[string]schema.Attribute{
		"file": schema.StringAttribute{Optional: true},
	}}

	t.Run("null_value", func(t *testing.T) {
		req := planmodifier.StringRequest{
			Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
				"file": tftypes.NewValue(tftypes.String, nil),
			}),
		}
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.PlanValue.IsNull() {
			t.Fatal("expected plan value to be null")
		}
	})

	t.Run("unknown_value", func(t *testing.T) {
		req := planmodifier.StringRequest{
			Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
				"file": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
		}
		resp := planmodifier.StringResponse{PlanValue: types.StringValue("keep")}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.PlanValue.IsNull() {
			t.Fatal("expected plan value to be null")
		}
	})

	t.Run("checksum_error", func(t *testing.T) {
		tmpDir := t.TempDir()
		req := planmodifier.StringRequest{
			Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
				"file": tftypes.NewValue(tftypes.String, tmpDir),
			}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics to have error")
		}
	})

	t.Run("empty_checksum", func(t *testing.T) {
		nonexistent := filepath.Join(t.TempDir(), "missing.txt")
		req := planmodifier.StringRequest{
			Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
				"file": tftypes.NewValue(tftypes.String, nonexistent),
			}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if !resp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics to have error")
		}
		if !resp.PlanValue.IsUnknown() {
			t.Fatal("expected plan value to be unknown")
		}
	})

	t.Run("checksum_success", func(t *testing.T) {
		tmpDir := t.TempDir()
		file := filepath.Join(tmpDir, "file.txt")
		if err := os.WriteFile(file, []byte("content"), 0600); err != nil {
			t.Fatalf("write file: %v", err)
		}
		checksum, err := helpers.CalculateSHA256(file)
		if err != nil {
			t.Fatalf("checksum: %v", err)
		}

		req := planmodifier.StringRequest{
			Plan: newPlan(t, schemaDef, map[string]tftypes.Value{
				"file": tftypes.NewValue(tftypes.String, file),
			}),
		}
		resp := planmodifier.StringResponse{}

		modifier.PlanModifyString(ctx, req, &resp)

		if resp.PlanValue.ValueString() != checksum {
			t.Fatalf("expected checksum %q, got %q", checksum, resp.PlanValue.ValueString())
		}
	})
}
