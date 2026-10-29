// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitCheckContextTimeout_NoTimeout(t *testing.T) {
	ctx := context.Background()
	err := helpers.CheckContextTimeout(ctx, "test operation")
	if err != nil {
		t.Errorf("Expected no error for non-cancelled context, got: %v", err)
	}
}

func TestUnitCheckContextTimeout_WithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := helpers.CheckContextTimeout(ctx, "test operation")
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}

	expectedErrorSubstring := "timed out during test operation"
	if !containsString(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedErrorSubstring, err.Error())
	}
}

// readTimeoutRaw builds a raw config/state value of the form `timeouts = { read = "<read>" }`.
func readTimeoutRaw(ctx context.Context, read string) tftypes.Value {
	timeoutsType := timeouts.Attributes(ctx, timeouts.Opts{Read: true}).GetType().TerraformType(ctx)
	rootType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"timeouts": timeoutsType}}

	return tftypes.NewValue(rootType, map[string]tftypes.Value{
		"timeouts": tftypes.NewValue(timeoutsType, map[string]tftypes.Value{
			"read": tftypes.NewValue(tftypes.String, read),
		}),
	})
}

func TestUnitEnterRequestContext_ResourceRead_AppliesReadTimeout(t *testing.T) {
	ctx := context.Background()
	req := resource.ReadRequest{
		State: tfsdk.State{
			Schema: resourceschema.Schema{
				Attributes: map[string]resourceschema.Attribute{
					"timeouts": timeouts.Attributes(ctx, timeouts.Opts{Read: true}),
				},
			},
			Raw: readTimeoutRaw(ctx, "2m"),
		},
	}

	ctx, exitContext := helpers.EnterRequestContext(ctx, helpers.TypeInfo{TypeName: "test"}, req)
	defer exitContext()

	deadline, ok := ctx.Deadline()
	require.True(t, ok, "expected resource read context to have a deadline")
	require.WithinDuration(t, time.Now().Add(2*time.Minute), deadline, 5*time.Second)
}

// Reproduces https://github.com/microsoft/terraform-provider-power-platform/issues/1247
func TestUnitEnterRequestContext_DataSourceRead_AppliesReadTimeout(t *testing.T) {
	ctx := context.Background()
	// Mirrors how data sources declare timeouts in this provider.
	req := datasource.ReadRequest{
		Config: tfsdk.Config{
			Schema: datasourceschema.Schema{
				Attributes: map[string]datasourceschema.Attribute{
					"timeouts": timeouts.Attributes(ctx, timeouts.Opts{Read: true}),
				},
			},
			Raw: readTimeoutRaw(ctx, "2m"),
		},
	}

	ctx, exitContext := helpers.EnterRequestContext(ctx, helpers.TypeInfo{TypeName: "test"}, req)
	defer exitContext()

	deadline, ok := ctx.Deadline()
	require.True(t, ok, "expected data source read context to have a deadline")
	require.WithinDuration(t, time.Now().Add(2*time.Minute), deadline, 5*time.Second)
}

// Helper function to check if a string contains a substring.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()))
}
