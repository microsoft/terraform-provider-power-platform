// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type TokenExpiredError struct {
	Message string
}

func (e *TokenExpiredError) Error() string {
	return e.Message
}

type Auth struct {
	config *config.ProviderConfig
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
		config: configValue,
	}
}

func (client *Auth) AuthenticateClientCertificate(ctx context.Context, scopes []string) (string, time.Time, error) {

	cert, key, err := helpers.ConvertBase64ToCert(client.config.ClientCertificateRaw, client.config.ClientCertificatePassword)
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	azureCertCredentials, err := azidentity.NewClientCertificateCredential(
		client.config.TenantId,
		client.config.ClientId,
		cert,
		key,
		&azidentity.ClientCertificateCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: client.config.Cloud,
			},
		},
	)
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}
	accessToken, err := azureCertCredentials.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}
	return accessToken.Token, accessToken.ExpiresOn, nil
}

//nolint:unused-receiver
func (client *Auth) AuthenticateUsingCli(ctx context.Context, scopes []string) (string, time.Time, error) {
	azureCLICredentials, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	accessToken, err := azureCLICredentials.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateClientSecret(ctx context.Context, scopes []string) (string, time.Time, error) {
	clientSecretCredential, err := azidentity.NewClientSecretCredential(
		client.config.TenantId,
		client.config.ClientId,
		client.config.ClientSecret, &azidentity.ClientSecretCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: client.config.Cloud,
			},
		})
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	accessToken, err := clientSecretCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: client.config.TenantId,
	})

	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil

}

func NewOidcCredential(options *OidcCredentialOptions) (*OidcCredential, error) {
	c := &OidcCredential{
		requestToken:  options.RequestToken,
		requestUrl:    options.RequestUrl,
		token:         options.Token,
		tokenFilePath: options.TokenFilePath,
	}

	if c.requestToken == constants.EMPTY {
		return nil, fmt.Errorf("request Token is required for OIDC credential")
	}
	if c.requestUrl == constants.EMPTY {
		return nil, fmt.Errorf("request URL is required for OIDC credential")
	}
	if options.TenantID == constants.EMPTY {
		return nil, fmt.Errorf("tenant is required for OIDC credential")
	}
	if options.ClientID == constants.EMPTY {
		return nil, fmt.Errorf("client is required for OIDC credential")
	}

	cred, err := azidentity.NewClientAssertionCredential(options.TenantID, options.ClientID, c.getAssertion,
		&azidentity.ClientAssertionCredentialOptions{
			ClientOptions: options.ClientOptions,
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
	var creds []azcore.TokenCredential

	oidcCred, err := NewOidcCredential(&OidcCredentialOptions{
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

	if err == nil {
		creds = append(creds, oidcCred)
	} else {
		log.Printf("newDefaultAzureCredential failed to initialize oidc credential:\n\t%s", err.Error())
	}
	chain, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	accessToken, err := chain.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return constants.EMPTY, time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (w *OidcCredential) getAssertion(ctx context.Context) (string, error) {
	if w.token != constants.EMPTY {
		return w.token, nil
	}

	if w.tokenFilePath != constants.EMPTY {
		idTokenData, err := os.ReadFile(w.tokenFilePath)
		if err != nil {
			return constants.EMPTY, fmt.Errorf("reading token file: %v", err)
		}

		return string(idTokenData), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.requestUrl, http.NoBody)
	if err != nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: failed to build request")
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: cannot parse URL query")
	}

	if query.Get("audience") == constants.EMPTY {
		query.Set("audience", "api://AzureADTokenExchange")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.requestToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: cannot request token: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: cannot parse response: %v", err)
	}

	if statusCode := resp.StatusCode; statusCode < RESPONSE_200 || statusCode >= RESPONSE_300 {
		return constants.EMPTY, fmt.Errorf("getAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Count *int    `json:"count"`
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: cannot unmarshal response: %v", err)
	}

	if tokenRes.Value == nil {
		return constants.EMPTY, fmt.Errorf("getAssertion: nil JWT assertion received from OIDC provider")
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
	case client.config.IsOidcProvided():
		token, tokenExpiry, err = client.AuthenticateOIDC(ctx, scopes)
	case client.config.IsClientCertificateCredentialsProvided():
		token, tokenExpiry, err = client.AuthenticateClientCertificate(ctx, scopes)
	default:
		return nil, errors.New("no credentials provided")
	}

	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Token acquired (expire: %s): **********", tokenExpiry))
	return &token, err
}
