// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestUnitUUIDTypeStringAndValueType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	typ := UUIDType{}

	if typ.String() != "UUIDType" {
		t.Fatalf("unexpected type string: %s", typ.String())
	}

	if _, ok := typ.ValueType(ctx).(UUIDValue); !ok {
		t.Fatal("ValueType should return UUIDValue")
	}

	unknown := basetypes.NewStringUnknown()
	val, diags := typ.ValueFromString(ctx, unknown)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if uuidVal, ok := val.(UUIDValue); !ok || !uuidVal.IsUnknown() {
		t.Fatal("expected ValueFromString to return unknown UUIDValue")
	}
}

func TestUnitUUIDTypeEqualDifferentType(t *testing.T) {
	t.Parallel()

	typ := UUIDType{}
	if typ.Equal(basetypes.StringType{}) {
		t.Fatal("expected UUIDType to not equal StringType")
	}
}
