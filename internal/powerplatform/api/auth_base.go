package powerplatform_common

import (
	"context"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

var _ AuthInterface = &AuthImplementation{}

type TokeExpiredError struct {
	Message string
}

func (e *TokeExpiredError) Error() string {
	return e.Message
}

type AuthInterface interface {
	IsTokenExpiredOrEmpty() bool

	GetToken() (string, error)
	SetToken(string)

	SetTokenExpiry(time.Time)
	GetTokenExpiry() time.Time

	AuthClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error)
}

type AuthImplementation struct {
	Config      common.ProviderConfig
	Token       string
	TokenExpiry time.Time
}

type AuthBaseOperationInterface interface {
	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

func (client *AuthImplementation) AuthClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error) {
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

func (client *AuthImplementation) SetToken(token string) {
	client.Token = token
}

func (client *AuthImplementation) SetTokenExpiry(tokenExpiry time.Time) {
	client.TokenExpiry = tokenExpiry
}

func (client *AuthImplementation) GetTokenExpiry() time.Time {
	return client.TokenExpiry
}

func (client *AuthImplementation) GetToken() (string, error) {
	if client.IsTokenExpiredOrEmpty() {
		return "", &TokeExpiredError{"token is expired or empty"}
	}
	return client.Token, nil
}

func (client *AuthImplementation) IsTokenExpiredOrEmpty() bool {
	return client.Token == "" || time.Now().After(client.TokenExpiry)
}
