package powerplatform

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ PowerPlatformAuthInterface = &PowerPlatformAuthImplementation{}

type PowerPlatformAuthInterface interface {
	IsTokenExpiredOrEmpty() bool
	RefreshToken() (string, error)

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type PowerPlatformAuthImplementation struct {
	Token       string
	TokenExpiry time.Time
}

func (client *PowerPlatformAuthImplementation) IsTokenExpiredOrEmpty() bool {
	if client.Token == "" {
		return true
	} else {
		return time.Now().After(client.TokenExpiry)
	}
}

func (client *PowerPlatformAuthImplementation) RefreshToken() (string, error) {
	//todo implement token refresh
	panic("[RefreshToken] not implemented")
}

func (client *PowerPlatformAuthImplementation) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	//todo implement
	panic("[AuthenticateUserPass] not implemented")
}

func (client *PowerPlatformAuthImplementation) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	scopes := []string{"https://api.powerplatform.com/.default"}
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

func (client *PowerPlatformAuthImplementation) authClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error) {
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
