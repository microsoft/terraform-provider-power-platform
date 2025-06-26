// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitCheckContextTimeout_NoTimeout(t *testing.T) {
	ctx := context.Background()
	err := helpers.CheckContextTimeout(ctx, "test operation")
	if err != nil {
		t.Errorf("Expected no error for non-cancelled context, got: %v", err)
	}
}

func TestUnitCheckContextTimeout_WithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for the timeout to be reached
	time.Sleep(10 * time.Millisecond)

	err := helpers.CheckContextTimeout(ctx, "test operation")
	if err == nil {
		t.Error("Expected timeout error for expired context, got nil")
	}

	expectedErrorSubstring := "timed out during test operation"
	if !containsString(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedErrorSubstring, err.Error())
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
