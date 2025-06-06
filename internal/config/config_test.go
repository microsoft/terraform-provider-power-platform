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

func TestUnitProviderConfig_IsCliProvided(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		expected bool
	}{
		{
			name: "IsCliProvided returns true when UseCli is true",
			config: ProviderConfig{
				UseCli: true,
			},
			expected: true,
		},
		{
			name: "IsCliProvided returns false when UseCli is false",
			config: ProviderConfig{
				UseCli: false,
			},
			expected: false,
		},
		{
			name: "IsCliProvided returns false when UseCli is default (false)",
			config: ProviderConfig{
				UseDevCli: true,
				UseOidc:   true,
				UseMsi:    true,
				// UseCli not set, should default to false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsCliProvided()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitProviderConfig_IsOidcProvided(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		expected bool
	}{
		{
			name: "IsOidcProvided returns true when UseOidc is true",
			config: ProviderConfig{
				UseOidc: true,
			},
			expected: true,
		},
		{
			name: "IsOidcProvided returns false when UseOidc is false",
			config: ProviderConfig{
				UseOidc: false,
			},
			expected: false,
		},
		{
			name: "IsOidcProvided returns false when UseOidc is default (false)",
			config: ProviderConfig{
				UseCli:    true,
				UseDevCli: true,
				UseMsi:    true,
				// UseOidc not set, should default to false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsOidcProvided()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitProviderConfig_IsSystemManagedIdentityProvided(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		expected bool
	}{
		{
			name: "IsSystemManagedIdentityProvided returns true when UseMsi is true and ClientId is empty",
			config: ProviderConfig{
				UseMsi:   true,
				ClientId: "",
			},
			expected: true,
		},
		{
			name: "IsSystemManagedIdentityProvided returns false when UseMsi is false",
			config: ProviderConfig{
				UseMsi:   false,
				ClientId: "",
			},
			expected: false,
		},
		{
			name: "IsSystemManagedIdentityProvided returns false when UseMsi is true but ClientId is not empty",
			config: ProviderConfig{
				UseMsi:   true,
				ClientId: "some-client-id",
			},
			expected: false,
		},
		{
			name: "IsSystemManagedIdentityProvided returns false when UseMsi is default (false)",
			config: ProviderConfig{
				UseCli:    true,
				UseDevCli: true,
				UseOidc:   true,
				// UseMsi not set, should default to false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsSystemManagedIdentityProvided()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitProviderConfig_IsUserManagedIdentityProvided(t *testing.T) {
	tests := []struct {
		name     string
		config   ProviderConfig
		expected bool
	}{
		{
			name: "IsUserManagedIdentityProvided returns true when UseMsi is true and ClientId is not empty",
			config: ProviderConfig{
				UseMsi:   true,
				ClientId: "some-client-id",
			},
			expected: true,
		},
		{
			name: "IsUserManagedIdentityProvided returns false when UseMsi is false",
			config: ProviderConfig{
				UseMsi:   false,
				ClientId: "some-client-id",
			},
			expected: false,
		},
		{
			name: "IsUserManagedIdentityProvided returns false when UseMsi is true but ClientId is empty",
			config: ProviderConfig{
				UseMsi:   true,
				ClientId: "",
			},
			expected: false,
		},
		{
			name: "IsUserManagedIdentityProvided returns false when UseMsi is default (false)",
			config: ProviderConfig{
				UseCli:    true,
				UseDevCli: true,
				UseOidc:   true,
				// UseMsi not set, should default to false
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsUserManagedIdentityProvided()
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

func TestUnitProviderConfig_AuthMethods_Independence(t *testing.T) {
	// Test that all authentication methods work independently
	config := ProviderConfig{
		UseCli:    true,
		UseDevCli: true,
		UseOidc:   true,
		UseMsi:    true,
		ClientId:  "test-client-id",
	}

	assert.True(t, config.IsCliProvided(), "IsCliProvided should be true")
	assert.True(t, config.IsDevCliProvided(), "IsDevCliProvided should be true")
	assert.True(t, config.IsOidcProvided(), "IsOidcProvided should be true")
	assert.True(t, config.IsUserManagedIdentityProvided(), "IsUserManagedIdentityProvided should be true (UseMsi=true, ClientId=not empty)")
	assert.False(t, config.IsSystemManagedIdentityProvided(), "IsSystemManagedIdentityProvided should be false (UseMsi=true, ClientId=not empty)")
}
