// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

type mockTokenCredential struct {
	getTokenCallCount int32
	tokenValue        string
	expiresOn         time.Time
	err               error
}

func (m *mockTokenCredential) GetToken(ctx context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	atomic.AddInt32(&m.getTokenCallCount, 1)
	if m.err != nil {
		return azcore.AccessToken{}, m.err
	}
	return azcore.AccessToken{
		Token:     m.tokenValue,
		ExpiresOn: m.expiresOn,
	}, nil
}

// generateTestPFX creates a self-signed RSA certificate, packages it into a password-protected
// PKCS#12 (PFX) blob, and returns the base64-encoded PFX data for use in tests.
func generateTestPFX(t *testing.T, password string) string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	pfxBytes, err := pkcs12.Modern.Encode(key, cert, nil, password)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(pfxBytes)
}

func writeCliStub(t *testing.T, dir, name, output string) string {
	t.Helper()
	filename := name
	script := "#!/bin/sh\necho '" + output + "'\n"
	perm := os.FileMode(0700)

	if runtime.GOOS == "windows" {
		filename = name + ".cmd"
		script = "@echo off\r\necho " + output + "\r\n"
		perm = 0600
	}

	path := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(path, []byte(script), perm))
	if runtime.GOOS != "windows" {
		require.NoError(t, os.Chmod(path, 0700))
	}
	return path
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
			wg.Go(func() {
				cred, err := authClient.getOrCreateCredential(ctx, credType, factory)
				results <- result{credType: credType, cred: cred, err: err}
			})
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
	ta.credentialsMutex.Lock()
	defer ta.credentialsMutex.Unlock()
	holder := &credentialHolder{credential: cred}
	holder.once.Do(func() {})
	ta.credentials[credType] = holder
}

func (ta *testableAuth) injectCredentialError(credType credentialType, err error) {
	ta.credentialsMutex.Lock()
	defer ta.credentialsMutex.Unlock()
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

func TestUnit_CliTokenCache_CachesAndReusesTokens(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	// Set a token in the cache
	cacheKey := "cli:https://management.azure.com/.default"
	expiresOn := time.Now().Add(1 * time.Hour)
	authClient.setCachedCliToken(cacheKey, "cached-token", expiresOn)

	// Retrieve the token
	token, retrievedExpiry, found := authClient.getCachedCliToken(cacheKey)

	assert.True(t, found, "Token should be found in cache")
	assert.Equal(t, "cached-token", token)
	assert.Equal(t, expiresOn, retrievedExpiry)
}

func TestUnit_CliTokenCache_ReturnsNotFoundForMissingKey(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	token, expiresOn, found := authClient.getCachedCliToken("nonexistent-key")

	assert.False(t, found, "Token should not be found")
	assert.Empty(t, token)
	assert.True(t, expiresOn.IsZero())
}

func TestUnit_CliTokenCache_ReturnsNotFoundForExpiredToken(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	// Set a token that expires in less than 5 minutes (within the buffer)
	cacheKey := "cli:https://management.azure.com/.default"
	expiresOn := time.Now().Add(3 * time.Minute)
	authClient.setCachedCliToken(cacheKey, "expiring-token", expiresOn)

	// Token should not be returned because it's within the 5-minute buffer
	token, _, found := authClient.getCachedCliToken(cacheKey)

	assert.False(t, found, "Token expiring soon should not be returned")
	assert.Empty(t, token)
}

func TestUnit_CliTokenCache_ReturnsValidToken(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	// Set a token that expires in more than 5 minutes
	cacheKey := "cli:https://management.azure.com/.default"
	expiresOn := time.Now().Add(30 * time.Minute)
	authClient.setCachedCliToken(cacheKey, "valid-token", expiresOn)

	// Token should be returned
	token, _, found := authClient.getCachedCliToken(cacheKey)

	assert.True(t, found, "Valid token should be returned")
	assert.Equal(t, "valid-token", token)
}

func TestUnit_CliTokenCache_DifferentScopesAreSeparate(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	expiresOn := time.Now().Add(1 * time.Hour)

	// Set tokens for different scopes
	authClient.setCachedCliToken("cli:scope1", "token1", expiresOn)
	authClient.setCachedCliToken("cli:scope2", "token2", expiresOn)

	// Retrieve and verify separate tokens
	token1, _, found1 := authClient.getCachedCliToken("cli:scope1")
	token2, _, found2 := authClient.getCachedCliToken("cli:scope2")

	assert.True(t, found1)
	assert.True(t, found2)
	assert.Equal(t, "token1", token1)
	assert.Equal(t, "token2", token2)
}

func TestUnit_CliTokenCache_ConcurrentAccess(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	const numGoroutines = 50
	var wg sync.WaitGroup

	expiresOn := time.Now().Add(1 * time.Hour)

	// Writers
	for i := range numGoroutines {
		wg.Go(func() {
			cacheKey := fmt.Sprintf("cli:scope%d", i%5)
			authClient.setCachedCliToken(cacheKey, fmt.Sprintf("token%d", i), expiresOn)
		})
	}

	// Readers
	for range numGoroutines {
		wg.Go(func() {
			authClient.getCachedCliToken("cli:scope0")
		})
	}

	wg.Wait()

	// Should not panic or deadlock - verify cache still works
	_, _, found := authClient.getCachedCliToken("cli:scope0")
	assert.True(t, found, "Token should be found after concurrent access")
}

func TestUnit_NewAuthBase_InitializesCliTokenCache(t *testing.T) {
	providerConfig := &config.ProviderConfig{}
	authClient := NewAuthBase(providerConfig)

	assert.NotNil(t, authClient.cliTokens, "CLI token cache should be initialized")
	assert.Empty(t, authClient.cliTokens, "CLI token cache should be empty initially")
}

func TestUnitTokenExpiredError_Error(t *testing.T) {
	t.Parallel()

	err := (&TokenExpiredError{Message: "expired"}).Error()
	require.Equal(t, "expired", err)
}

func TestUnitNewOidcCredential_ValidationErrors(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{}
	auth := NewAuthBase(cfg)

	_, err := auth.NewOidcCredential(&OidcCredentialOptions{})
	require.Error(t, err)

	_, err = auth.NewOidcCredential(&OidcCredentialOptions{RequestToken: "token"})
	require.Error(t, err)

	_, err = auth.NewOidcCredential(&OidcCredentialOptions{RequestToken: "token", RequestUrl: "https://example.com"})
	require.Error(t, err)

	_, err = auth.NewOidcCredential(&OidcCredentialOptions{RequestToken: "token", RequestUrl: "https://example.com", TenantID: "tid"})
	require.Error(t, err)
}

func TestUnitOidcCredential_GetAssertion(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cred := &OidcCredential{token: "direct-token"}
	token, err := cred.getAssertion(ctx)
	require.NoError(t, err)
	require.Equal(t, "direct-token", token)

	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "token")
	require.NoError(t, os.WriteFile(filePath, []byte("file-token"), 0600))
	cred = &OidcCredential{tokenFilePath: filePath}
	token, err = cred.getAssertion(ctx)
	require.NoError(t, err)
	require.Equal(t, "file-token", token)

	successSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer request-token", r.Header.Get("Authorization"))
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"value":"http-token"}`))
	}))
	defer successSrv.Close()

	cred = &OidcCredential{requestToken: "request-token", requestUrl: successSrv.URL}
	token, err = cred.getAssertion(ctx)
	require.NoError(t, err)
	require.Equal(t, "http-token", token)

	errorSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad"}`))
	}))
	defer errorSrv.Close()

	cred = &OidcCredential{requestToken: "request-token", requestUrl: errorSrv.URL}
	_, err = cred.getAssertion(ctx)
	require.Error(t, err)

	missingValueSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"value":null}`))
	}))
	defer missingValueSrv.Close()

	cred = &OidcCredential{requestToken: "request-token", requestUrl: missingValueSrv.URL}
	_, err = cred.getAssertion(ctx)
	require.Error(t, err)
}

func TestUnitGetTokenForScopes_Paths(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{TestMode: true}
	auth := NewAuthBase(cfg)

	token, err := auth.GetTokenForScopes(context.Background(), []string{"scope"})
	require.NoError(t, err)
	require.Equal(t, "test_mode_mock_token_value", *token)

	cfg.TestMode = false
	auth = NewAuthBase(cfg)
	_, err = auth.GetTokenForScopes(context.Background(), []string{"scope"})
	require.EqualError(t, err, "no credentials provided")
}

func TestUnitAuthenticateClientCertificate_InvalidCert(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{ClientCertificateRaw: "not-base64", ClientCertificatePassword: "pass"}
	auth := NewAuthBase(cfg)
	_, _, err := auth.AuthenticateClientCertificate(context.Background(), []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateClientSecret_MissingTenant(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{ClientId: "cid", ClientSecret: "secret"}
	auth := NewAuthBase(cfg)
	_, _, err := auth.AuthenticateClientSecret(context.Background(), []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateCliVariants(t *testing.T) {
	t.Parallel()
	cfg := &config.ProviderConfig{}
	auth := NewAuthBase(cfg)

	_, _, err := auth.AuthenticateUsingCli(context.Background(), []string{"scope"})
	require.Error(t, err)

	_, _, err = auth.AuthenticateUsingAzureDeveloperCli(context.Background(), []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateManagedIdentity_ContextCancelled(t *testing.T) {
	cfg := &config.ProviderConfig{ClientId: "client"}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateUserManagedIdentity(ctx, []string{"scope"})
	require.Error(t, err)

	_, _, err = auth.AuthenticateSystemManagedIdentity(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateAzDOWorkloadIdentityFederation_Validation(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{}
	auth := NewAuthBase(cfg)

	_, _, err := auth.AuthenticateAzDOWorkloadIdentityFederation(context.Background(), []string{"scope"})
	require.EqualError(t, err, "tenant ID must be provided to use Azure DevOps Workload Identity Federation")

	cfg.TenantId = "tenant"
	_, _, err = auth.AuthenticateAzDOWorkloadIdentityFederation(context.Background(), []string{"scope"})
	require.EqualError(t, err, "client ID must be provided to use Azure DevOps Workload Identity Federation")

	cfg.ClientId = "client"
	_, _, err = auth.AuthenticateAzDOWorkloadIdentityFederation(context.Background(), []string{"scope"})
	require.EqualError(t, err, "the Azure DevOps service connection ID could not be found")

	cfg.AzDOServiceConnectionID = "service"
	_, _, err = auth.AuthenticateAzDOWorkloadIdentityFederation(context.Background(), []string{"scope"})
	require.EqualError(t, err, "could not obtain an OIDC request token for Azure DevOps Workload Identity Federation")
}

func TestUnitAuthenticateOIDC_ContextCancelled(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{
		TenantId:         "tenant",
		ClientId:         "client",
		OidcRequestToken: "request-token",
		OidcRequestUrl:   "https://example.com/oidc",
		OidcToken:        "oidc-token",
	}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateOIDC(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateClientCertificate_ContextCancelled(t *testing.T) {
	t.Parallel()

	pfx := generateTestPFX(t, "pass")
	cfg := &config.ProviderConfig{ClientCertificateRaw: pfx, ClientCertificatePassword: "pass"}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateClientCertificate(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateClientSecret_ContextCancelled(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{TenantId: "tenant", ClientId: "client", ClientSecret: "secret"}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateClientSecret(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateCliVariants_WithStub(t *testing.T) {
	// Not parallel because the test modifies PATH to use stub CLI executables.
	tmp := t.TempDir()
	writeCliStub(t, tmp, "az", `{"accessToken":"token","expiresOn":"2099-01-01 00:00:00.000000"}`)
	writeCliStub(t, tmp, "azd", `{"token":"token","expiresOn":"2099-01-01T00:00:00Z"}`)
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfg := &config.ProviderConfig{}
	auth := NewAuthBase(cfg)
	token, _, err := auth.AuthenticateUsingCli(context.Background(), []string{"scope"})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token, _, err = auth.AuthenticateUsingAzureDeveloperCli(context.Background(), []string{"scope"})
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestUnitAuthenticateManagedIdentity_Attempt(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{UseMsi: true, ClientId: "client"}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateUserManagedIdentity(ctx, []string{"scope"})
	require.Error(t, err)

	cfg.ClientId = ""
	_, _, err = auth.AuthenticateSystemManagedIdentity(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitAuthenticateAzDOWorkloadIdentityFederation_Cancelled(t *testing.T) {
	t.Parallel()

	cfg := &config.ProviderConfig{
		TenantId:                "tenant",
		ClientId:                "client",
		AzDOServiceConnectionID: "service",
		OidcRequestToken:        "token",
	}
	auth := NewAuthBase(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := auth.AuthenticateAzDOWorkloadIdentityFederation(ctx, []string{"scope"})
	require.Error(t, err)
}

func TestUnitGetTokenForScopes_AllBranches(t *testing.T) {
	pfx := generateTestPFX(t, "pass")
	tmp := t.TempDir()
	writeCliStub(t, tmp, "az", `{"accessToken":"token","expiresOn":"2099-01-01 00:00:00.000000"}`)
	writeCliStub(t, tmp, "azd", `{"token":"token","expiresOn":"2099-01-01T00:00:00Z"}`)
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	tests := []config.ProviderConfig{
		{TenantId: "tenant", ClientId: "client", ClientSecret: "secret"},
		{UseCli: true},
		{UseDevCli: true},
		{UseOidc: true, AzDOServiceConnectionID: "service", TenantId: "tenant", ClientId: "client", OidcRequestToken: "token"},
		{UseOidc: true, TenantId: "tenant", ClientId: "client", OidcRequestToken: "token", OidcRequestUrl: "https://example.com", OidcToken: "token"},
		{ClientCertificateRaw: pfx, ClientCertificatePassword: "pass"},
		{UseMsi: true, ClientId: "client"},
		{UseMsi: true},
	}

	for _, cfg := range tests {
		cfgCopy := cfg
		auth := NewAuthBase(&cfgCopy)
		ctx := context.Background()
		var cancel context.CancelFunc
		if cfgCopy.UseMsi {
			ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		}
		token, err := auth.GetTokenForScopes(ctx, []string{"scope"})
		if cancel != nil {
			cancel()
		}
		if err == nil {
			require.NotNil(t, token)
		}
	}
}
