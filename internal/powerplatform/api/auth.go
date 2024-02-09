package powerplatform

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
)

type TokeExpiredError struct {
	Message string
}

func (e *TokenExpiredError) Error() string {
	return e.Message
}

type Auth struct {
	config *config.ProviderConfig
}

func NewAuthBase(config *config.ProviderConfig) *Auth {
	return &Auth{
		config: config,
	}
}

func (client *Auth) GetAuthority(tenantid string) string {
	return constants.OAUTH_AUTHORITY_URL + tenantid
}

func (client *Auth) AuthenticateUsingCli(ctx context.Context, scopes []string) (string, time.Time, error) {
	azureCLICredentials, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := azureCLICredentials.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil
}

func (client *Auth) AuthenticateClientSecret(ctx context.Context, scopes []string) (string, time.Time, error) {
	clientSecretCredential, err := azidentity.NewClientSecretCredential(
		client.config.Credentials.TenantId,
		client.config.Credentials.ClientId,
		client.config.Credentials.ClientSecret, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	accessToken, err := clientSecretCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes:   scopes,
		TenantID: client.config.Credentials.TenantId,
	})

	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken.Token, accessToken.ExpiresOn, nil

}

func (client *Auth) GetTokenForScopes(ctx context.Context, scopes []string) (*string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[GetTokenForScope] Getting token for scope: '%s'", strings.Join(scopes, ",")))

	token := ""
	tokenExpiry := time.Time{}
	var err error

	switch {
	case client.config.Credentials.IsClientSecretCredentialsProvided():
		token, tokenExpiry, err = client.AuthenticateClientSecret(ctx, scopes)
	case client.config.Credentials.IsCliProvided():
		token, tokenExpiry, err = client.AuthenticateUsingCli(ctx, scopes)
	default:
		return nil, errors.New("no credentials provided")
	}

	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, fmt.Sprintf("Token acquired (expire: %s): **********", tokenExpiry))
	return &token, err
}
