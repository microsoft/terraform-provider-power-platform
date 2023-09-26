package powerplatform_api_bapi

import (
	"context"
	"errors"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ BapiAuthInterface = &BapiAuth{}

type BapiAuthInterface interface {
	GetBase() api.AuthInterface

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type BapiAuth struct {
	BaseAuth api.AuthInterface
}

func (client *BapiAuth) GetBase() api.AuthInterface {
	return client.BaseAuth
}

func (client *BapiAuth) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
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

	client.BaseAuth.SetToken(authResult.AccessToken)
	client.BaseAuth.SetTokenExpiry(authResult.ExpiresOn)

	return client.BaseAuth.GetToken()
}

func (client *BapiAuth) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	scopes := []string{"https://service.powerapps.com//.default"}
	token, expiry, err := client.BaseAuth.AuthClientSecret(ctx, scopes, tenantId, applicationId, secret)
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
		}
		return "", err
	}

	client.BaseAuth.SetToken(token)
	client.BaseAuth.SetTokenExpiry(expiry)

	return client.BaseAuth.GetToken()
}
