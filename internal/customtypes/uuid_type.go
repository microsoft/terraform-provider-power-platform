// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*UUIDType)(nil)

type UUIDType struct {
	basetypes.StringType
}

func (t UUIDType) Equal(o attr.Type) bool {
	// Support both value and pointer types for comparison
	switch v := o.(type) {
	case UUIDType:
		return t.StringType.Equal(v.StringType)
	case *UUIDType:
		return t.StringType.Equal(v.StringType)
	default:
		return false
	}
}

func (t UUIDType) String() string {
	return "UUIDType"
}

func (t UUIDType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := UUIDValue{
		StringValue: in,
	}

	return value, nil
}

func (t UUIDType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t UUIDType) ValueType(_ context.Context) attr.Value {
	return UUIDValue{}
}
