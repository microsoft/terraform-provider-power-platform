// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewUUIDNull() UUID {
	return UUID{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewUUIDUnknown() UUID {
	return UUID{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewUUIDValue(value string) UUID {
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewUUIDPointerValue(value *string) UUID {
	if value == nil {
		return NewUUIDNull()
	}

	return NewUUIDValue(*value)
}

func NewUUIDValueMust(value string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(value).ValueUUID()
}

func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(*value).ValueUUID()
}
