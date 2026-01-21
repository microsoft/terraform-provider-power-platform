// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type TokenExpiredError struct {
	Message string
}

func (e *TokenExpiredError) Error() string {
	return e.Message
}

type credentialType string

const (
	credTypeClientSecret          credentialType = "client_secret"
	credTypeClientCertificate     credentialType = "client_certificate"
	credTypeCLI                   credentialType = "cli"
	credTypeDevCLI                credentialType = "dev_cli"
	credTypeAzDOPipelines         credentialType = "azdo_pipelines"
	credTypeOIDC                  credentialType = "oidc"
	credTypeUserManagedIdentity   credentialType = "user_managed_identity"
	credTypeSystemManagedIdentity credentialType = "system_managed_identity"
)

type credentialHolder struct {
	credential azcore.TokenCredential
	once       sync.Once
	err        error
}

type cachedToken struct {
	token     string
	expiresOn time.Time
}

type Auth struct {
	config *config.ProviderConfig

	credentials      map[credentialType]*credentialHolder
	credentialsMutex sync.RWMutex

	cliTokens      map[string]*cachedToken
	cliTokensMutex sync.RWMutex
}

type OidcCredential struct {
	requestToken  string
	requestUrl    string
	token         string
	tokenFilePath string
	cred          *azidentity.ClientAssertionCredential
}

type OidcCredentialOptions struct {
	azcore.ClientOptions
	TenantID      string
	ClientID      string
	RequestToken  string
	RequestUrl    string
	Token         string
	TokenFilePath string
}

func NewAuthBase(configValue *config.ProviderConfig) *Auth {
	return &Auth{
		config:      configValue,
		credentials: make(map[credentialType]*credentialHolder),
		cliTokens:   make(map[string]*cachedToken),
	}
}

func (client *Auth) getOrCreateCredential(ctx context.Context, credType credentialType, factory func() (azcore.TokenCredential, error)) (azcore.TokenCredential, error) {
	client.credentialsMutex.RLock()
	holder, exists := client.credentials[credType]
	client.credentialsMutex.RUnlock()

	if !exists {
		client.credentialsMutex.Lock()
		holder, exists = client.credentials[credType]
		if !exists {
			holder = &credentialHolder{}
			client.credentials[credType] = holder
			tflog.Debug(ctx, fmt.Sprintf("Created credential holder for type: %s", credType))
		}
		client.credentialsMutex.Unlock()
	}

	holder.once.Do(func() {
		tflog.Debug(ctx, fmt.Sprintf("Initializing credential for type: %s", credType))
		holder.credential, holder.err = factory()
		if holder.err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to create credential for type %s: %s", credType, holder.err.Error()))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Successfully created credential for type: %s", credType))
		}
	})

	return holder.credential, holder.err
}

func (client *Auth) getCachedCliToken(cacheKey string) (string, time.Time, bool) {
	client.cliTokensMutex.RLock()
	defer client.cliTokensMutex.RUnlock()

	cached, exists := client.cliTokens[cacheKey]
	if !exists {
		return "", time.Time{}, false
	}

	if time.Now().Add(5 * time.Minute).Before(cached.expiresOn) {
		return cached.token, cached.expiresOn, true
	}

	return "", time.Time{}, false
}

func (client *Auth) setCachedCliToken(cacheKey string, token string, expiresOn time.Time) {
	client.cliTokensMutex.Lock()
	defer client.cliTokensMutex.Unlock()

	client.cliTokens[cacheKey] = &cachedToken{
		token:     token,
		expiresOn: expiresOn,
	}
}

func (client *Auth) AuthenticateClientCertificate(ctx context.Context, scopes []string) (string, time.Time, error) {
	cred, err := client.getOrCreateCredential(ctx, credTypeClientCertificate, func() (azcore.TokenCredential, error) {
		cert, key, certErr := helpers.ConvertBase64ToCert(client.config.ClientCertificateRaw, client.config.ClientCertificatePassword)
		if certErr != nil {
			return nil, certErr
		}

		return azidentity.NewClientCertificateCredential(
			client.config.TenantId,
			client.config.ClientId,
			cert,
			key,
			&azidentity.ClientCertificateCredentialOptions{
				AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
				ClientOptions: azcore.ClientOptions{
					Cloud: client.config.Cloud,
				},
			},
		)
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}
	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateUsingCli(ctx context.Context, scopes []string) (string, time.Time, error) {
	cacheKey := "cli:" + strings.Join(scopes, ",")
	if token, expiresOn, found := client.getCachedCliToken(cacheKey); found {
		tflog.Debug(ctx, "Using cached token for Azure CLI credential")
		return token, expiresOn, nil
	}

	cred, err := client.getOrCreateCredential(ctx, credTypeCLI, func() (azcore.TokenCredential, error) {
		return azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{
			AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
			TenantID:                   client.config.TenantId,
		})
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	client.setCachedCliToken(cacheKey, accessToken.Token, accessToken.ExpiresOn)
	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateUsingAzureDeveloperCli(ctx context.Context, scopes []string) (string, time.Time, error) {
	cacheKey := "devcli:" + strings.Join(scopes, ",")
	if token, expiresOn, found := client.getCachedCliToken(cacheKey); found {
		tflog.Debug(ctx, "Using cached token for Azure Developer CLI credential")
		return token, expiresOn, nil
	}

	cred, err := client.getOrCreateCredential(ctx, credTypeDevCLI, func() (azcore.TokenCredential, error) {
		return azidentity.NewAzureDeveloperCLICredential(&azidentity.AzureDeveloperCLICredentialOptions{
			AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
			TenantID:                   client.config.TenantId,
		})
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	client.setCachedCliToken(cacheKey, accessToken.Token, accessToken.ExpiresOn)
	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateClientSecret(ctx context.Context, scopes []string) (string, time.Time, error) {
	cred, err := client.getOrCreateCredential(ctx, credTypeClientSecret, func() (azcore.TokenCredential, error) {
		return azidentity.NewClientSecretCredential(
			client.config.TenantId,
			client.config.ClientId,
			client.config.ClientSecret,
			&azidentity.ClientSecretCredentialOptions{
				AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
				ClientOptions: azcore.ClientOptions{
					Cloud: client.config.Cloud,
				},
			},
		)
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}
	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) NewOidcCredential(options *OidcCredentialOptions) (*OidcCredential, error) {
	c := &OidcCredential{
		requestToken:  options.RequestToken,
		requestUrl:    options.RequestUrl,
		token:         options.Token,
		tokenFilePath: options.TokenFilePath,
	}

	if c.requestToken == "" {
		return nil, errors.New("request Token is required for OIDC credential")
	}
	if c.requestUrl == "" {
		return nil, errors.New("request URL is required for OIDC credential")
	}
	if options.TenantID == "" {
		return nil, errors.New("tenant is required for OIDC credential")
	}
	if options.ClientID == "" {
		return nil, errors.New("client is required for OIDC credential")
	}

	cred, err := azidentity.NewClientAssertionCredential(options.TenantID, options.ClientID, c.getAssertion,
		&azidentity.ClientAssertionCredentialOptions{
			AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
			ClientOptions:              options.ClientOptions,
		})
	if err != nil {
		return nil, err
	}

	c.cred = cred
	return c, nil
}

func (w *OidcCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return w.cred.GetToken(ctx, opts)
}

func (client *Auth) AuthenticateOIDC(ctx context.Context, scopes []string) (string, time.Time, error) {
	cred, err := client.getOrCreateCredential(ctx, credTypeOIDC, func() (azcore.TokenCredential, error) {
		oidcCred, oidcErr := client.NewOidcCredential(&OidcCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: client.config.Cloud,
			},
			TenantID:      client.config.TenantId,
			ClientID:      client.config.ClientId,
			RequestToken:  client.config.OidcRequestToken,
			RequestUrl:    client.config.OidcRequestUrl,
			Token:         client.config.OidcToken,
			TokenFilePath: client.config.OidcTokenFilePath,
		})
		if oidcErr != nil {
			return nil, oidcErr
		}

		return azidentity.NewChainedTokenCredential([]azcore.TokenCredential{oidcCred}, nil)
	})
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("newDefaultAzureCredential failed to initialize oidc credential:\n\t%s", err.Error()))
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateUserManagedIdentity(ctx context.Context, scopes []string) (string, time.Time, error) {
	cred, err := client.getOrCreateCredential(ctx, credTypeUserManagedIdentity, func() (azcore.TokenCredential, error) {
		return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(client.config.ClientId),
			ClientOptions: azcore.ClientOptions{
				Cloud: client.config.Cloud,
			},
		})
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateSystemManagedIdentity(ctx context.Context, scopes []string) (string, time.Time, error) {
	cred, err := client.getOrCreateCredential(ctx, credTypeSystemManagedIdentity, func() (azcore.TokenCredential, error) {
		return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: client.config.Cloud,
			},
		})
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateAzDOWorkloadIdentityFederation(ctx context.Context, scopes []string) (string, time.Time, error) {
	if client.config.TenantId == "" {
		return "", time.Time{}, errors.New("tenant ID must be provided to use Azure DevOps Workload Identity Federation")
	}
	if client.config.ClientId == "" {
		return "", time.Time{}, errors.New("client ID must be provided to use Azure DevOps Workload Identity Federation")
	}
	if client.config.AzDOServiceConnectionID == "" {
		return "", time.Time{}, errors.New("the Azure DevOps service connection ID could not be found")
	}
	if client.config.OidcRequestToken == "" {
		return "", time.Time{}, errors.New("could not obtain an OIDC request token for Azure DevOps Workload Identity Federation")
	}

	cred, err := client.getOrCreateCredential(ctx, credTypeAzDOPipelines, func() (azcore.TokenCredential, error) {
		return azidentity.NewAzurePipelinesCredential(
			client.config.TenantId,
			client.config.ClientId,
			client.config.AzDOServiceConnectionID,
			client.config.OidcRequestToken,
			&azidentity.AzurePipelinesCredentialOptions{
				AdditionallyAllowedTenants: client.config.AuxiliaryTenantIDs,
				ClientOptions: azcore.ClientOptions{
					Cloud: client.config.Cloud,
				},
			},
		)
	})
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := cred.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

// createTokenRequestOptions creates TokenRequestOptions with CAE support when enabled.
func (client *Auth) createTokenRequestOptions(ctx context.Context, scopes []string) policy.TokenRequestOptions {
	tokenOptions := policy.TokenRequestOptions{
		Scopes: scopes,
	}

	// Add TenantID for ClientSecret authentication if it's available
	if client.config.TenantId != "" {
		tokenOptions.TenantID = client.config.TenantId
	}

	// Enable CAE if configured
	if client.config.EnableContinuousAccessEvaluation {
		tokenOptions.EnableCAE = true
		tflog.Debug(ctx, "Continuous Access Evaluation (CAE) is enabled for token requests")
	}

	return tokenOptions
}

func (w *OidcCredential) getAssertion(ctx context.Context) (string, error) {
	if w.token != "" {
		return w.token, nil
	}

	if w.tokenFilePath != "" {
		idTokenData, err := os.ReadFile(w.tokenFilePath)
		if err != nil {
			return "", fmt.Errorf("reading token file: %w", err)
		}

		return string(idTokenData), nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", w.requestUrl, http.NoBody)
	if err != nil {
		return "", errors.New("getAssertion: failed to build request")
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return "", errors.New("getAssertion: cannot parse URL query")
	}

	if query.Get("audience") == "" {
		query.Set("audience", "api://AzureADTokenExchange")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.requestToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot parse response: %w", err)
	}

	if statusCode := resp.StatusCode; statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("getAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Count *int    `json:"count"`
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return "", fmt.Errorf("getAssertion: cannot unmarshal response: %w", err)
	}

	if tokenRes.Value == nil {
		return "", errors.New("getAssertion: nil JWT assertion received from OIDC provider")
	}

	return *tokenRes.Value, nil
}

func (client *Auth) GetTokenForScopes(ctx context.Context, scopes []string) (*string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[GetTokenForScope] Getting token for scope: '%s'", strings.Join(scopes, ",")))

	if client.config.TestMode {
		token := "test_mode_mock_token_value"
		return &token, nil
	}

	token := ""
	var tokenExpiry time.Time
	var err error

	switch {
	case client.config.IsClientSecretCredentialsProvided():
		token, tokenExpiry, err = client.AuthenticateClientSecret(ctx, scopes)
	case client.config.IsCliProvided():
		token, tokenExpiry, err = client.AuthenticateUsingCli(ctx, scopes)
	case client.config.IsDevCliProvided():
		token, tokenExpiry, err = client.AuthenticateUsingAzureDeveloperCli(ctx, scopes)
	case client.config.IsAzDOWorkloadIdentityFederationProvided():
		token, tokenExpiry, err = client.AuthenticateAzDOWorkloadIdentityFederation(ctx, scopes)
	case client.config.IsOidcProvided():
		token, tokenExpiry, err = client.AuthenticateOIDC(ctx, scopes)
	case client.config.IsClientCertificateCredentialsProvided():
		token, tokenExpiry, err = client.AuthenticateClientCertificate(ctx, scopes)
	case client.config.IsUserManagedIdentityProvided():
		token, tokenExpiry, err = client.AuthenticateUserManagedIdentity(ctx, scopes)
	case client.config.IsSystemManagedIdentityProvided():
		token, tokenExpiry, err = client.AuthenticateSystemManagedIdentity(ctx, scopes)
	default:
		return nil, errors.New("no credentials provided")
	}

	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Token acquired (expire: %s): **********", tokenExpiry))
	return &token, err
}
