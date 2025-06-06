// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitProviderConfig_IsDevCliProvided(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		expected bool
	}{
		{
			name: "IsDevCliProvided returns true when UseDevCli is true",
			config: ProviderConfig{
				UseDevCli: true,
			},
			expected: true,
		},
		{
			name: "IsDevCliProvided returns false when UseDevCli is false",
			config: ProviderConfig{
				UseDevCli: false,
			},
			expected: false,
		},
		{
			name: "IsDevCliProvided returns false when UseDevCli is default (false)",
			config: ProviderConfig{
				UseCli:  true,
				UseOidc: true,
				UseMsi:  true,
				// UseDevCli not set, should default to false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsDevCliProvided()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitProviderConfig_IsDevCliProvided_vs_OtherMethods(t *testing.T) {
	// Ensure IsDevCliProvided works independently of other auth methods
	config := ProviderConfig{
		UseCli:    true,
		UseDevCli: true,
		UseOidc:   true,
		UseMsi:    true,
	}

	assert.True(t, config.IsCliProvided(), "IsCliProvided should be true")
	assert.True(t, config.IsDevCliProvided(), "IsDevCliProvided should be true")
	assert.True(t, config.IsOidcProvided(), "IsOidcProvided should be true")
	assert.True(t, config.IsSystemManagedIdentityProvided(), "IsSystemManagedIdentityProvided should be true (UseMsi=true, ClientId=empty)")
}
