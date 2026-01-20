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
)

func TestUnitUUIDValueEqualAndType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	val := NewUUIDValue(validUUID)

	if _, ok := val.Type(ctx).(UUIDType); !ok {
		t.Fatal("expected Type to return UUIDType")
	}

	if !val.Equal(NewUUIDValue(validUUID)) {
		t.Fatal("expected UUIDValue equality for matching values")
	}

	if val.Equal(basetypes.NewStringValue(validUUID)) {
		t.Fatal("expected UUIDValue equality to return false for other types")
	}
}

func TestUnitUUIDValueSemanticEqualsWrongType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	val := NewUUIDValue(validUUID)
	other := basetypes.NewStringValue(validUUID)

	equal, diags := val.StringSemanticEquals(ctx, other)
	if !diags.HasError() {
		t.Fatal("expected diagnostics for unexpected value type")
	}
	if equal {
		t.Fatal("expected equality to be false when diagnostics are present")
	}
}

func TestUnitUUIDValueValidateNullUnknown(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := xattr.ValidateAttributeRequest{Path: path.Root("id")}

	nullValue := NewUUIDNull()
	nullResp := &xattr.ValidateAttributeResponse{Diagnostics: diag.Diagnostics{}}
	nullValue.ValidateAttribute(ctx, request, nullResp)
	if nullResp.Diagnostics.HasError() {
		t.Fatal("expected no diagnostics for null value")
	}

	unknownValue := NewUUIDUnknown()
	unknownResp := &xattr.ValidateAttributeResponse{Diagnostics: diag.Diagnostics{}}
	unknownValue.ValidateAttribute(ctx, request, unknownResp)
	if unknownResp.Diagnostics.HasError() {
		t.Fatal("expected no diagnostics for unknown value")
	}

	paramResp := &function.ValidateParameterResponse{}
	nullValue.ValidateParameter(ctx, function.ValidateParameterRequest{Position: 0}, paramResp)
	if paramResp.Error != nil {
		t.Fatal("expected no error for null parameter")
	}

	paramResp = &function.ValidateParameterResponse{}
	unknownValue.ValidateParameter(ctx, function.ValidateParameterRequest{Position: 0}, paramResp)
	if paramResp.Error != nil {
		t.Fatal("expected no error for unknown parameter")
	}
}
