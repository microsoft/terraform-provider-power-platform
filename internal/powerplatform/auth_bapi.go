package powerplatform

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ BapiAuthInterface = &BapiAuthImplementation{}

type TokeExpiredError struct {
	message string
}

func (e *TokeExpiredError) Error() string {
	return e.message
}

type BapiAuthInterface interface {
	IsTokenExpiredOrEmpty() bool
	GetCurrentToken() (string, error)

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type BapiAuthImplementation struct {
	Config      ProviderConfig
	Token       string
	TokenExpiry time.Time
}

func (client *BapiAuthImplementation) GetCurrentToken() (string, error) {
	if client.IsTokenExpiredOrEmpty() {
		return "", &TokeExpiredError{"token is expired or empty"}
	}
	return client.Token, nil
}

func (client *BapiAuthImplementation) IsTokenExpiredOrEmpty() bool {
	return client.Token == "" || time.Now().After(client.TokenExpiry)
}

func (client *BapiAuthImplementation) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	scopes := []string{"https://service.powerapps.com//.default"}
	publicClientApplicationID := "1950a258-227b-4e31-a9cf-717495945fc2"
	authority := "https://login.microsoftonline.com/" + tenantId

	publicClientApp, err := public.New(publicClientApplicationID, public.WithAuthority(authority))
	if err != nil {
		return "", err
	}

	authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, username, password)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.Token = authResult.AccessToken
	client.TokenExpiry = authResult.ExpiresOn

	return client.Token, nil
}

func (client *BapiAuthImplementation) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	scopes := []string{"https://service.powerapps.com//.default"}
	token, expiry, err := client.authClientSecret(ctx, scopes, tenantId, applicationId, secret)
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
		}
		return "", err
	}
	client.Token = token
	client.TokenExpiry = expiry
	return client.Token, nil
}

func (client *BapiAuthImplementation) authClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error) {
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
