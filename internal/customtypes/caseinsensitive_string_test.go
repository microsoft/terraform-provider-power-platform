// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestUnitCaseInsensitiveString_SemanticEquals(t *testing.T) {
	cases := []struct {
		name  string
		old   string
		new   string
		equal bool
	}{
		{"identical", "Power Platform Reader", "Power Platform Reader", true},
		{"case only", "Power Platform Reader", "power platform READER", true},
		{"different", "Power Platform Reader", "Power Platform Administrator", false},
		{"whitespace matters", "Power Platform Reader", "Power  Platform Reader", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			old := NewCaseInsensitiveStringValue(tc.old)
			equal, diags := old.StringSemanticEquals(context.Background(), NewCaseInsensitiveStringValue(tc.new))
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if equal != tc.equal {
				t.Errorf("expected semantic equality %v for %q vs %q, got %v", tc.equal, tc.old, tc.new, equal)
			}
		})
	}
}

func TestUnitCaseInsensitiveString_RejectsForeignType(t *testing.T) {
	v := NewCaseInsensitiveStringValue("x")
	_, diags := v.StringSemanticEquals(context.Background(), basetypes.NewStringValue("x"))
	if !diags.HasError() {
		t.Error("a foreign value type must produce a diagnostic")
	}
}
