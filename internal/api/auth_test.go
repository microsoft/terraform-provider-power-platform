// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTokenCredential struct {
	getTokenCallCount int32
	tokenValue        string
	expiresOn         time.Time
	err               error
}

func (m *mockTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	atomic.AddInt32(&m.getTokenCallCount, 1)
	if m.err != nil {
		return azcore.AccessToken{}, m.err
	}
	return azcore.AccessToken{
		Token:     m.tokenValue,
		ExpiresOn: m.expiresOn,
	}, nil
}

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

func TestUnit_GetOrCreateCredential_ReusesCredential(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	factoryCallCount := 0
	mockCred := &mockTokenCredential{tokenValue: "test-token"}

	factory := func() (azcore.TokenCredential, error) {
		factoryCallCount++
		return mockCred, nil
	}

	// First call - should create the credential
	cred1, err1 := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
	require.NoError(t, err1)
	assert.NotNil(t, cred1)
	assert.Equal(t, 1, factoryCallCount, "Factory should be called exactly once on first call")

	// Second call - should reuse the credential
	cred2, err2 := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
	require.NoError(t, err2)
	assert.NotNil(t, cred2)
	assert.Equal(t, 1, factoryCallCount, "Factory should not be called again on second call")

	// Third call - should still reuse the credential
	cred3, err3 := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
	require.NoError(t, err3)
	assert.NotNil(t, cred3)
	assert.Equal(t, 1, factoryCallCount, "Factory should not be called again on third call")

	// Verify same instance is returned
	assert.Same(t, cred1, cred2, "Same credential instance should be returned")
	assert.Same(t, cred2, cred3, "Same credential instance should be returned")
}

func TestUnit_GetOrCreateCredential_DifferentTypesAreSeparate(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	clientSecretCallCount := 0
	cliCallCount := 0

	clientSecretCred := &mockTokenCredential{tokenValue: "client-secret-token"}
	cliCred := &mockTokenCredential{tokenValue: "cli-token"}

	clientSecretFactory := func() (azcore.TokenCredential, error) {
		clientSecretCallCount++
		return clientSecretCred, nil
	}

	cliFactory := func() (azcore.TokenCredential, error) {
		cliCallCount++
		return cliCred, nil
	}

	// Create client secret credential
	cred1, err := authClient.getOrCreateCredential(ctx, credTypeClientSecret, clientSecretFactory)
	require.NoError(t, err)
	assert.Equal(t, 1, clientSecretCallCount)
	assert.Equal(t, 0, cliCallCount)

	// Create CLI credential
	cred2, err := authClient.getOrCreateCredential(ctx, credTypeCLI, cliFactory)
	require.NoError(t, err)
	assert.Equal(t, 1, clientSecretCallCount)
	assert.Equal(t, 1, cliCallCount)

	// Verify different instances
	assert.NotSame(t, cred1, cred2, "Different credential types should return different instances")
	assert.Same(t, cred1, clientSecretCred)
	assert.Same(t, cred2, cliCred)

	// Re-request both - should reuse
	cred1Again, _ := authClient.getOrCreateCredential(ctx, credTypeClientSecret, clientSecretFactory)
	cred2Again, _ := authClient.getOrCreateCredential(ctx, credTypeCLI, cliFactory)

	assert.Equal(t, 1, clientSecretCallCount, "Factory should not be called again")
	assert.Equal(t, 1, cliCallCount, "Factory should not be called again")
	assert.Same(t, cred1, cred1Again)
	assert.Same(t, cred2, cred2Again)
}

func TestUnit_GetOrCreateCredential_CachesErrors(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	factoryCallCount := 0
	expectedErr := errors.New("authentication failed")

	factory := func() (azcore.TokenCredential, error) {
		factoryCallCount++
		return nil, expectedErr
	}

	// First call - should get the error
	cred1, err1 := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
	assert.Nil(t, cred1)
	assert.Equal(t, expectedErr, err1)
	assert.Equal(t, 1, factoryCallCount, "Factory should be called exactly once")

	// Second call - should return the cached error without calling factory again
	cred2, err2 := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
	assert.Nil(t, cred2)
	assert.Equal(t, expectedErr, err2)
	assert.Equal(t, 1, factoryCallCount, "Factory should not be called again")
}

func TestUnit_GetOrCreateCredential_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	var factoryCallCount int32
	mockCred := &mockTokenCredential{tokenValue: "concurrent-token"}

	factory := func() (azcore.TokenCredential, error) {
		atomic.AddInt32(&factoryCallCount, 1)
		time.Sleep(10 * time.Millisecond)
		return mockCred, nil
	}

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	credentials := make([]azcore.TokenCredential, numGoroutines)
	errs := make([]error, numGoroutines)

	for i := range numGoroutines {
		go func(index int) {
			defer wg.Done()
			cred, err := authClient.getOrCreateCredential(ctx, credTypeClientSecret, factory)
			credentials[index] = cred
			errs[index] = err
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(1), atomic.LoadInt32(&factoryCallCount),
		"Factory should be called exactly once even with concurrent access")

	for i, cred := range credentials {
		assert.NoError(t, errs[i])
		assert.Same(t, mockCred, cred, "All goroutines should receive the same credential instance")
	}
}

func TestUnit_GetOrCreateCredential_ConcurrentAccessDifferentTypes(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	var clientSecretCallCount, cliCallCount, devCliCallCount int32
	clientSecretCred := &mockTokenCredential{tokenValue: "client-secret"}
	cliCred := &mockTokenCredential{tokenValue: "cli"}
	devCliCred := &mockTokenCredential{tokenValue: "dev-cli"}

	factories := map[credentialType]func() (azcore.TokenCredential, error){
		credTypeClientSecret: func() (azcore.TokenCredential, error) {
			atomic.AddInt32(&clientSecretCallCount, 1)
			time.Sleep(5 * time.Millisecond)
			return clientSecretCred, nil
		},
		credTypeCLI: func() (azcore.TokenCredential, error) {
			atomic.AddInt32(&cliCallCount, 1)
			time.Sleep(5 * time.Millisecond)
			return cliCred, nil
		},
		credTypeDevCLI: func() (azcore.TokenCredential, error) {
			atomic.AddInt32(&devCliCallCount, 1)
			time.Sleep(5 * time.Millisecond)
			return devCliCred, nil
		},
	}

	expectedCreds := map[credentialType]azcore.TokenCredential{
		credTypeClientSecret: clientSecretCred,
		credTypeCLI:          cliCred,
		credTypeDevCLI:       devCliCred,
	}

	const numGoroutinesPerType = 20
	var wg sync.WaitGroup

	type result struct {
		credType credentialType
		cred     azcore.TokenCredential
		err      error
	}
	results := make(chan result, numGoroutinesPerType*3)

	for credType, factory := range factories {
		for range numGoroutinesPerType {
			wg.Add(1)
			go func(ct credentialType, f func() (azcore.TokenCredential, error)) {
				defer wg.Done()
				cred, err := authClient.getOrCreateCredential(ctx, ct, f)
				results <- result{credType: ct, cred: cred, err: err}
			}(credType, factory)
		}
	}

	wg.Wait()
	close(results)

	assert.Equal(t, int32(1), atomic.LoadInt32(&clientSecretCallCount))
	assert.Equal(t, int32(1), atomic.LoadInt32(&cliCallCount))
	assert.Equal(t, int32(1), atomic.LoadInt32(&devCliCallCount))

	for r := range results {
		assert.NoError(t, r.err)
		assert.Same(t, expectedCreds[r.credType], r.cred,
			"Credential for type %s should be the expected instance", r.credType)
	}
}

func TestUnit_NewAuthBase_InitializesCredentialsMap(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	assert.NotNil(t, authClient.credentials, "Credentials map should be initialized")
	assert.Empty(t, authClient.credentials, "Credentials map should be empty initially")
	assert.Equal(t, providerConfig, authClient.config, "Config should be set correctly")
}

func TestUnit_CredentialTypes_AreUnique(t *testing.T) {
	types := []credentialType{
		credTypeClientSecret,
		credTypeClientCertificate,
		credTypeCLI,
		credTypeDevCLI,
		credTypeAzDOPipelines,
		credTypeOIDC,
		credTypeUserManagedIdentity,
		credTypeSystemManagedIdentity,
	}

	seen := make(map[credentialType]bool)
	for _, ct := range types {
		assert.False(t, seen[ct], "Credential type %s should be unique", ct)
		seen[ct] = true
	}
}

func TestUnit_GetTokenForScopes_TestMode(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{
		TestMode: true,
	}
	authClient := NewAuthBase(providerConfig)

	scopes := []string{"https://management.azure.com/.default"}

	token, err := authClient.GetTokenForScopes(ctx, scopes)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "test_mode_mock_token_value", *token)

	assert.Empty(t, authClient.credentials, "No credentials should be created in test mode")
}

func TestUnit_GetTokenForScopes_NoCredentials(t *testing.T) {
	ctx := context.Background()
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	scopes := []string{"https://management.azure.com/.default"}

	token, err := authClient.GetTokenForScopes(ctx, scopes)

	assert.Nil(t, token)
	assert.Error(t, err)
	assert.Equal(t, "no credentials provided", err.Error())
}

func TestUnit_TokenExpiredError(t *testing.T) {
	err := &TokenExpiredError{Message: "token has expired"}
	assert.Equal(t, "token has expired", err.Error())
}

type testableAuth struct {
	*Auth
}

func (ta *testableAuth) injectCredential(credType credentialType, cred azcore.TokenCredential) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	holder := &credentialHolder{credential: cred}
	holder.once.Do(func() {})
	ta.credentials[credType] = holder
}

func (ta *testableAuth) injectCredentialError(credType credentialType, err error) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	holder := &credentialHolder{err: err}
	holder.once.Do(func() {})
	ta.credentials[credType] = holder
}

func TestUnit_AuthenticateUserManagedIdentity_WithMockedCredential(t *testing.T) {
	ctx := context.Background()
	expiresOn := time.Now().Add(1 * time.Hour)
	mockCred := &mockTokenCredential{
		tokenValue: "user-managed-identity-token",
		expiresOn:  expiresOn,
	}

	providerConfig := &config.ProviderConfig{
		UseMsi:   true,
		ClientId: "test-client-id",
	}
	authClient := &testableAuth{Auth: NewAuthBase(providerConfig)}
	authClient.injectCredential(credTypeUserManagedIdentity, mockCred)

	scopes := []string{"https://management.azure.com/.default"}

	token, tokenExpiry, err := authClient.AuthenticateUserManagedIdentity(ctx, scopes)

	require.NoError(t, err)
	assert.Equal(t, "user-managed-identity-token", token)
	assert.Equal(t, expiresOn, tokenExpiry)
	assert.Equal(t, int32(1), atomic.LoadInt32(&mockCred.getTokenCallCount))
}

func TestUnit_AuthenticateUserManagedIdentity_GetTokenError(t *testing.T) {
	ctx := context.Background()
	mockCred := &mockTokenCredential{
		err: errors.New("failed to get token for user managed identity"),
	}

	providerConfig := &config.ProviderConfig{
		UseMsi:   true,
		ClientId: "test-client-id",
	}
	authClient := &testableAuth{Auth: NewAuthBase(providerConfig)}
	authClient.injectCredential(credTypeUserManagedIdentity, mockCred)

	scopes := []string{"https://management.azure.com/.default"}

	token, tokenExpiry, err := authClient.AuthenticateUserManagedIdentity(ctx, scopes)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, time.Time{}, tokenExpiry)
	assert.Contains(t, err.Error(), "failed to get token")
}

func TestUnit_AuthenticateUserManagedIdentity_CredentialCreationError(t *testing.T) {
	ctx := context.Background()

	providerConfig := &config.ProviderConfig{
		UseMsi:   true,
		ClientId: "test-client-id",
	}
	authClient := &testableAuth{Auth: NewAuthBase(providerConfig)}
	authClient.injectCredentialError(credTypeUserManagedIdentity, errors.New("credential creation failed"))

	scopes := []string{"https://management.azure.com/.default"}

	token, tokenExpiry, err := authClient.AuthenticateUserManagedIdentity(ctx, scopes)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, time.Time{}, tokenExpiry)
	assert.Equal(t, "credential creation failed", err.Error())
}

func TestUnit_AuthenticateSystemManagedIdentity_WithMockedCredential(t *testing.T) {
	ctx := context.Background()
	expiresOn := time.Now().Add(1 * time.Hour)
	mockCred := &mockTokenCredential{
		tokenValue: "system-managed-identity-token",
		expiresOn:  expiresOn,
	}

	providerConfig := &config.ProviderConfig{
		UseMsi: true,
	}
	authClient := &testableAuth{Auth: NewAuthBase(providerConfig)}
	authClient.injectCredential(credTypeSystemManagedIdentity, mockCred)

	scopes := []string{"https://management.azure.com/.default"}

	token, tokenExpiry, err := authClient.AuthenticateSystemManagedIdentity(ctx, scopes)

	require.NoError(t, err)
	assert.Equal(t, "system-managed-identity-token", token)
	assert.Equal(t, expiresOn, tokenExpiry)
	assert.Equal(t, int32(1), atomic.LoadInt32(&mockCred.getTokenCallCount))
}

func TestUnit_AuthenticateSystemManagedIdentity_GetTokenError(t *testing.T) {
	ctx := context.Background()
	mockCred := &mockTokenCredential{
		err: errors.New("system identity token retrieval failed"),
	}

	providerConfig := &config.ProviderConfig{
		UseMsi: true,
	}
	authClient := &testableAuth{Auth: NewAuthBase(providerConfig)}
	authClient.injectCredential(credTypeSystemManagedIdentity, mockCred)

	scopes := []string{"https://management.azure.com/.default"}

	token, tokenExpiry, err := authClient.AuthenticateSystemManagedIdentity(ctx, scopes)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, time.Time{}, tokenExpiry)
}
