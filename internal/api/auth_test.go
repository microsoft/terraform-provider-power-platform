// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestUnitCreateTokenRequestOptions(t *testing.T) {
	scopes := []string{"https://management.azure.com/.default"}
	ctx := context.Background()

	testCases := []struct {
		name               string
		tenantID           string
		enableCAE          bool
		expectedScopes     []string
		expectedTenantID   string
		expectedCAEEnabled bool
	}{
		{
			name:               "Default without CAE or TenantID",
			tenantID:           "",
			enableCAE:          false,
			expectedScopes:     scopes,
			expectedTenantID:   "",
			expectedCAEEnabled: false,
		},
		{
			name:               "Only CAE Enabled",
			tenantID:           "",
			enableCAE:          true,
			expectedScopes:     scopes,
			expectedTenantID:   "",
			expectedCAEEnabled: true,
		},
		{
			name:               "Only TenantID Provided",
			tenantID:           "11111111-1111-1111-1111-111111111111",
			enableCAE:          false,
			expectedScopes:     scopes,
			expectedTenantID:   "11111111-1111-1111-1111-111111111111",
			expectedCAEEnabled: false,
		},
		{
			name:               "Both CAE and TenantID Provided",
			tenantID:           "22222222-2222-2222-2222-222222222222",
			enableCAE:          true,
			expectedScopes:     scopes,
			expectedTenantID:   "22222222-2222-2222-2222-222222222222",
			expectedCAEEnabled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			providerConfig := &config.ProviderConfig{
				TenantId:                         tc.tenantID,
				EnableContinuousAccessEvaluation: tc.enableCAE,
			}
			authClient := NewAuthBase(providerConfig)

			// Act - Since createTokenRequestOptions is private, access it directly for testing
			tokenOptions := authClient.createTokenRequestOptions(ctx, scopes)

			// Assert
			assert.Equal(t, tc.expectedScopes, tokenOptions.Scopes, "Scopes should match")
			assert.Equal(t, tc.expectedTenantID, tokenOptions.TenantID, "TenantID should match")
			assert.Equal(t, tc.expectedCAEEnabled, tokenOptions.EnableCAE, "EnableCAE should match")
		})
	}
}

func TestUnitAuthenticateUsingAzureDeveloperCli_ConfigurationCheck(t *testing.T) {
	// Test that the configuration is set up properly for dev CLI authentication
	testCases := []struct {
		name           string
		useDevCli      bool
		expectedMethod string
	}{
		{
			name:           "UseDevCli true should trigger dev CLI auth",
			useDevCli:      true,
			expectedMethod: "dev_cli",
		},
		{
			name:           "UseDevCli false should not trigger dev CLI auth",
			useDevCli:      false,
			expectedMethod: "none",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			providerConfig := &config.ProviderConfig{
				UseDevCli: tc.useDevCli,
			}

			// Act & Assert
			if tc.expectedMethod == "dev_cli" {
				assert.True(t, providerConfig.IsDevCliProvided(), "IsDevCliProvided should return true")
			} else {
				assert.False(t, providerConfig.IsDevCliProvided(), "IsDevCliProvided should return false")
			}
		})
	}
}
