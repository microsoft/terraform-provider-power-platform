package powerplatform

import (
	"context"
	"errors"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

var _ AuthBaseOperationInterface = &PowerPlatformAuth{}

type PowerPlatformAuth struct {
	baseAuth *AuthBase
}

func NewPowerPlatformAuth(baseAuth *AuthBase) *PowerPlatformAuth {
	return &PowerPlatformAuth{
		baseAuth: baseAuth,
	}
}

func (client *PowerPlatformAuth) AuthenticateUserPass(ctx context.Context, credentials *powerplatform_common.ProviderCredentials) (string, error) {
	scopes := []string{"https://api.powerplatform.com/.default"}
	publicClientApplicationID := "1950a258-227b-4e31-a9cf-717495945fc2"
	authority := "https://login.microsoftonline.com/" + credentials.TenantId

	publicClientApp, err := public.New(publicClientApplicationID, public.WithAuthority(authority))
	if err != nil {
		return "", err
	}

	authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.baseAuth.SetToken(authResult.AccessToken)
	client.baseAuth.SetTokenExpiry(authResult.ExpiresOn)

	return client.baseAuth.GetToken()
}

func (client *PowerPlatformAuth) AuthenticateClientSecret(ctx context.Context, credentials *powerplatform_common.ProviderCredentials) (string, error) {
	scopes := []string{"https://api.powerplatform.com/.default"}
	token, expiry, err := client.baseAuth.AuthClientSecret(ctx, scopes, credentials)
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
		}
		return "", err
	}
	client.baseAuth.SetToken(token)
	client.baseAuth.SetTokenExpiry(expiry)
	return client.baseAuth.GetToken()
}
