package powerplatform

import (
	"context"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	powerplatform_common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type TokeExpiredError struct {
	Message string
}

func (e *TokeExpiredError) Error() string {
	return e.Message
}

type AuthBase struct {
	config      *powerplatform_common.ProviderConfig
	token       string
	tokenExpiry time.Time
}

func NewAuthBase(config *powerplatform_common.ProviderConfig) *AuthBase {
	return &AuthBase{
		config: config,
	}
}

type AuthBaseOperationInterface interface {
	AuthenticateUserPass(ctx context.Context, credentials *powerplatform_common.ProviderCredentials) (string, error)
	AuthenticateClientSecret(ctx context.Context, credentials *powerplatform_common.ProviderCredentials) (string, error)
}

func (client *AuthBase) AuthClientSecret(ctx context.Context, scopes []string, credentials *powerplatform_common.ProviderCredentials) (string, time.Time, error) {
	authority := "https://login.microsoftonline.com/" + credentials.TenantId

	cred, err := confidential.NewCredFromSecret(credentials.Secret)
	if err != nil {
		return "", time.Time{}, err
	}

	confidentialClientApp, err := confidential.New(authority, credentials.ClientId, cred)
	if err != nil {
		return "", time.Time{}, err
	}

	authResult, err := confidentialClientApp.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		return "", time.Time{}, err
	}

	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) SetToken(token string) {
	client.token = token
}

func (client *AuthBase) SetTokenExpiry(tokenExpiry time.Time) {
	client.tokenExpiry = tokenExpiry
}

func (client *AuthBase) GetTokenExpiry() time.Time {
	return client.tokenExpiry
}

func (client *AuthBase) GetToken() (string, error) {
	if client.IsTokenExpiredOrEmpty() {
		return "", &TokeExpiredError{"token is expired or empty"}
	}
	return client.token, nil
}

func (client *AuthBase) IsTokenExpiredOrEmpty() bool {
	return client.token == "" || time.Now().After(client.tokenExpiry)
}
