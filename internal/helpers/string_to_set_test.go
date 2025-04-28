// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestUnitStringSliceToSet_Success(t *testing.T) {
	// Test valid input
	slice := []string{"value1", "value2", "value3"}

	// Convert the slice to a set
	set, err := helpers.StringSliceToSet(slice)

	// Verify no error occurred
	assert.NoError(t, err)
	assert.False(t, set.IsNull())
	assert.False(t, set.IsUnknown())

	// Verify the set has the expected type and elements
	assert.Equal(t, types.StringType, set.ElementType(context.TODO()))

	// Convert set back to slice for comparison
	elements := helpers.SetToStringSlice(set)
	assert.ElementsMatch(t, slice, elements)
}

func TestUnitStringSliceToSet_EmptySlice(t *testing.T) {
	// Test with empty slice
	emptySlice := []string{}

	// Convert empty slice to set
	set, err := helpers.StringSliceToSet(emptySlice)

	// Verify no error for empty slice
	assert.NoError(t, err)
	assert.False(t, set.IsNull())
	assert.False(t, set.IsUnknown())

	// Verify the set is empty
	elements := helpers.SetToStringSlice(set)
	assert.Empty(t, elements)
}

func TestUnitStringSliceToSet_DuplicateValues(t *testing.T) {
	// Test with duplicate values - this is actually valid for a set and should work
	duplicateSlice := []string{"value1", "value1", "value2"}

	// Convert to set (duplicates will be removed)
	set, err := helpers.StringSliceToSet(duplicateSlice)

	// Verify no error
	assert.NoError(t, err)

	// Duplicate values are not removed in the set
	assert.Equal(t, 3, len(set.Elements()))
}
