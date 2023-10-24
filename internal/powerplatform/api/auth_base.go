package powerplatform

import (
	"context"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type TokeExpiredError struct {
	Message string
}

func (e *TokeExpiredError) Error() string {
	return e.Message
}

type AuthBase struct {
	config      *common.ProviderConfig
	token       string
	tokenExpiry time.Time
}

func NewAuthBase(config *common.ProviderConfig) *AuthBase {
	return &AuthBase{
		config: config,
	}
}

type AuthBaseOperationInterface interface {
	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

func (client *AuthBase) AuthClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error) {
	authority := "https://login.microsoftonline.com/" + tenantId

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	confidentialClientApp, err := confidential.New(authority, applicationId, cred)
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
