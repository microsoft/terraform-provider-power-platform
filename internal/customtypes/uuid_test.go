// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	validUUID   = "00000000-0000-0000-0000-000000000001"
	anotherUUID = "00000000-0000-0000-0000-000000000002"
)

func TestUnitUUIDConstructors(t *testing.T) {
	t.Parallel()

	if !NewUUIDNull().IsNull() {
		t.Fatalf("expected NewUUIDNull to be null")
	}

	if !NewUUIDUnknown().IsUnknown() {
		t.Fatalf("expected NewUUIDUnknown to be unknown")
	}

	v := NewUUIDValue(validUUID)
	if v.ValueString() != validUUID {
		t.Fatalf("expected value to match input")
	}

	if p := NewUUIDPointerValue(nil); !p.IsNull() {
		t.Fatalf("nil pointer should produce null UUID")
	}

	ptr := validUUID
	if p := NewUUIDPointerValue(&ptr); p.ValueString() != validUUID {
		t.Fatalf("pointer constructor did not propagate value")
	}
}

func TestUnitUUIDValueUUIDValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	valid, diags := NewUUIDValueMust(validUUID)
	if diags.HasError() {
		t.Fatalf("expected valid UUID diagnostics to be clean: %v", diags)
	}
	if valid.ValueString() != validUUID {
		t.Fatalf("unexpected value %s", valid.ValueString())
	}

	_, diags = NewUUIDValueMust("invalid")
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for invalid UUID")
	}

	_, diags = NewUUIDPointerValue(nil).ValueUUID()
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for nil pointer UUID")
	}

	nullValue := NewUUIDNull()
	if _, diags := nullValue.ValueUUID(); !diags.HasError() {
		t.Fatalf("expected diagnostics for null ValueUUID")
	}

	unknownValue := NewUUIDUnknown()
	if _, diags := unknownValue.ValueUUID(); !diags.HasError() {
		t.Fatalf("expected diagnostics for unknown ValueUUID")
	}

	// ValidateAttribute should flag invalid UUID strings.
	invalid := NewUUIDValue("not-a-uuid")
	attrResp := &xattr.ValidateAttributeResponse{Diagnostics: diag.Diagnostics{}}
	invalid.ValidateAttribute(ctx, xattr.ValidateAttributeRequest{Path: path.Root("id")}, attrResp)
	if !attrResp.Diagnostics.HasError() {
		t.Fatalf("expected attribute diagnostics for invalid uuid")
	}

	// ValidateParameter should also surface errors.
	paramResp := &function.ValidateParameterResponse{}
	invalid.ValidateParameter(ctx, function.ValidateParameterRequest{Position: 1}, paramResp)
	if paramResp.Error == nil {
		t.Fatalf("expected parameter validation error")
	}
}

func TestUnitUUIDSemanticEquality(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	oldValue := NewUUIDValue(validUUID)
	newValue := NewUUIDValue(validUUID)

	equal, diags := oldValue.StringSemanticEquals(ctx, newValue)
	if diags.HasError() || !equal {
		t.Fatalf("expected semantic equality with same UUID")
	}

	unequalValue := NewUUIDValue(anotherUUID)
	equal, diags = oldValue.StringSemanticEquals(ctx, unequalValue)
	if diags.HasError() || equal {
		t.Fatalf("expected semantic inequality for different UUIDs")
	}

	badValue := NewUUIDValue("not-a-uuid")
	equal, diags = badValue.StringSemanticEquals(ctx, unequalValue)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for invalid uuid comparison")
	}
	if equal {
		t.Fatalf("expected equality to be false when diagnostics are present")
	}
}

func TestUnitUUIDTypeConversions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	typ := UUIDType{}

	req := basetypes.NewStringValue(validUUID)
	val, diags := typ.ValueFromString(ctx, req)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if _, ok := val.(UUIDValue); !ok {
		t.Fatalf("ValueFromString should return UUIDValue")
	}

	tfVal := tftypes.NewValue(tftypes.String, validUUID)
	attrVal, err := typ.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		t.Fatalf("unexpected error converting from terraform value: %v", err)
	}

	if _, ok := attrVal.(UUIDValue); !ok {
		t.Fatalf("expected UUIDValue from terraform conversion")
	}

	badTf := tftypes.NewValue(tftypes.Bool, true)
	if _, err := typ.ValueFromTerraform(ctx, badTf); err == nil {
		t.Fatalf("expected error for invalid terraform type")
	}

	if !typ.Equal(UUIDType{}) {
		t.Fatalf("UUIDType should be equal to another UUIDType value")
	}

	if !typ.Equal(&UUIDType{}) {
		t.Fatalf("UUIDType should consider pointer receivers equal")
	}
}

func TestUnitUUIDPointerValueMustValid(t *testing.T) {
	t.Parallel()

	ptr := validUUID
	val, diags := NewUUIDPointerValueMust(&ptr)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}

	if val.ValueString() != validUUID {
		t.Fatalf("expected value to match input")
	}
}

func TestUnitUUIDPointerValueMustNilPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil pointer")
		}
	}()

	var ptr *string
	_, _ = NewUUIDPointerValueMust(ptr)
}
