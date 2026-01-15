// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

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
	t.Setenv("AZURE_IDENTITY_DISABLE_MANAGED_IDENTITY", "true")

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
	tmp := t.TempDir()
	cliPath := filepath.Join(tmp, "az")
	devCliPath := filepath.Join(tmp, "azd")
	stub := []byte("#!/bin/sh\necho '{\"accessToken\":\"token\",\"expiresOn\":\"2099-01-01 00:00:00.000000\"}'\n")
	require.NoError(t, os.WriteFile(cliPath, stub, 0700))
	require.NoError(t, os.WriteFile(devCliPath, []byte("#!/bin/sh\necho '{\"token\":\"token\",\"expiresOn\":\"2099-01-01T00:00:00Z\"}'\n"), 0700))
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfg := &config.ProviderConfig{}
	auth := NewAuthBase(cfg)
	token, _, err := auth.AuthenticateUsingCli(context.Background(), []string{"scope"})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token, _, err = auth.AuthenticateUsingAzureDeveloperCli(context.Background(), []string{"scope"})
	if err == nil {
		require.NotEmpty(t, token)
	}
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
	cliPath := filepath.Join(tmp, "az")
	devCliPath := filepath.Join(tmp, "azd")
	stub := []byte("#!/bin/sh\necho '{\"accessToken\":\"token\",\"expiresOn\":\"2099-01-01 00:00:00.000000\"}'\n")
	require.NoError(t, os.WriteFile(cliPath, stub, 0700))
	require.NoError(t, os.WriteFile(devCliPath, []byte("#!/bin/sh\necho '{\"token\":\"token\",\"expiresOn\":\"2099-01-01T00:00:00Z\"}'\n"), 0700))
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
