// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitNormalizeNullableDescription_PreservesExplicitEmptyString(t *testing.T) {
	got := normalizeNullableDescription("", types.StringValue(""))
	if got.IsNull() {
		t.Fatal("expected explicit empty string to be preserved, got null")
	}
	if got.ValueString() != "" {
		t.Fatalf("expected empty string, got %q", got.ValueString())
	}
}

func TestUnitNormalizeNullableDescription_UsesNullWhenDescriptionMissing(t *testing.T) {
	got := normalizeNullableDescription("", types.StringNull())
	if !got.IsNull() {
		t.Fatalf("expected null for missing description, got %#v", got)
	}
}

func TestUnitNormalizeNullableDescription_PreservesNonEmptyValue(t *testing.T) {
	got := normalizeNullableDescription("Created by Terraform", types.StringNull())
	if got.IsNull() {
		t.Fatal("expected non-empty description, got null")
	}
	if got.ValueString() != "Created by Terraform" {
		t.Fatalf("expected non-empty description to be preserved, got %q", got.ValueString())
	}
}
