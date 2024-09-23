// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitTypeInfoString(t *testing.T) {
	t.Parallel()

	type testData struct {
		name     string
		typeInfo helpers.TypeInfo
		expected string
	}

	for _, testCase := range []testData{
		{
			name:     "simple",
			typeInfo: helpers.TypeInfo{ProviderTypeName: "provider", TypeName: "type"},
			expected: "provider_type",
		},
		{
			name:     "empty",
			typeInfo: helpers.TypeInfo{},
			expected: "powerplatform_", // Default provider name used
		},
		{
			name:     "empty provider",
			typeInfo: helpers.TypeInfo{TypeName: "type"},
			expected: "powerplatform_type", // Default provider name used
		},
		{
			name:     "empty type",
			typeInfo: helpers.TypeInfo{ProviderTypeName: "provider"},
			expected: "provider_",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.typeInfo.FullTypeName()
			if actual != testCase.expected {
				t.Errorf("expected %s, got %s", testCase.expected, actual)
			}
		})
	}
}
