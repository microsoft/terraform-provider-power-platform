// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = (*CaseInsensitiveStringType)(nil)
	_ basetypes.StringValuable                   = (*CaseInsensitiveStringValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*CaseInsensitiveStringValue)(nil)
)

// CaseInsensitiveStringType is a string whose values differing only by case are semantically
// equal. Use it for attributes the service matches case-insensitively, such as role definition
// display names. Semantic equality normalizes refresh and apply results against the prior value;
// suppressing a case-only CONFIG edit from planning a change additionally needs a plan modifier
// that folds a case-equal planned value back to state, since the framework does not consult
// semantic equality for configured value changes at plan time.
type CaseInsensitiveStringType struct {
	basetypes.StringType
}

func (t CaseInsensitiveStringType) Equal(o attr.Type) bool {
	switch v := o.(type) {
	case CaseInsensitiveStringType:
		return t.StringType.Equal(v.StringType)
	case *CaseInsensitiveStringType:
		return t.StringType.Equal(v.StringType)
	default:
		return false
	}
}

func (t CaseInsensitiveStringType) String() string {
	return "CaseInsensitiveStringType"
}

func (t CaseInsensitiveStringType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return CaseInsensitiveStringValue{StringValue: in}, nil
}

func (t CaseInsensitiveStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t CaseInsensitiveStringType) ValueType(_ context.Context) attr.Value {
	return CaseInsensitiveStringValue{}
}

type CaseInsensitiveString = CaseInsensitiveStringValue

type CaseInsensitiveStringValue struct {
	basetypes.StringValue
}

func (v CaseInsensitiveStringValue) Type(_ context.Context) attr.Type {
	return CaseInsensitiveStringType{}
}

func (v CaseInsensitiveStringValue) Equal(o attr.Value) bool {
	other, ok := o.(CaseInsensitiveStringValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v CaseInsensitiveStringValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(CaseInsensitiveStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return strings.EqualFold(v.ValueString(), newValue.ValueString()), diags
}

func NewCaseInsensitiveStringNull() CaseInsensitiveString {
	return CaseInsensitiveString{StringValue: basetypes.NewStringNull()}
}

func NewCaseInsensitiveStringValue(value string) CaseInsensitiveString {
	return CaseInsensitiveString{StringValue: basetypes.NewStringValue(value)}
}
