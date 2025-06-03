// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewUUIDNull returns a UUID representing a null value.
func NewUUIDNull() UUID {
	return UUID{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewUUIDUnknown returns a UUID representing an unknown value.
func NewUUIDUnknown() UUID {
	return UUID{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewUUIDValue returns a UUID initialized with the given string value.
func NewUUIDValue(value string) UUID {
	return UUID{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewUUIDPointerValue returns a UUID from a string pointer, or null if the pointer is nil.
func NewUUIDPointerValue(value *string) UUID {
	if value == nil {
		return NewUUIDNull()
	}

	return NewUUIDValue(*value)
}

// NewUUIDValueMust returns a UUID from a string value and validates it as a UUID.
func NewUUIDValueMust(value string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(value).ValueUUID()
}

// NewUUIDPointerValueMust returns a UUID from a string pointer and validates it as a UUID.
func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(*value).ValueUUID()
}
