package powerplatform

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
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

func (client *PowerPlatformAuth) AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthUsingCli(ctx, scopes, credentials)
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

func (client *PowerPlatformAuth) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthenticateUserPass(ctx, scopes, credentials)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.baseAuth.SetToken(token)
	client.baseAuth.SetTokenExpiry(expiry)

	return client.baseAuth.GetToken()
}

func (client *PowerPlatformAuth) AuthenticateClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
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
