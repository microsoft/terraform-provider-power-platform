package powerplatform

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ PowerPlatformAuthInterface = &PowerPlatformAuth{}

type PowerPlatformAuthInterface interface {
	GetBase() api.AuthInterface

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type PowerPlatformAuth struct {
	BaseAuth api.AuthInterface
}

func (client *PowerPlatformAuth) GetBase() api.AuthInterface {
	return client.BaseAuth
}

func (client *PowerPlatformAuth) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	//todo implement
	panic("[AuthenticateUserPass] not implemented")
}

func (client *PowerPlatformAuth) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	scopes := []string{"https://api.powerplatform.com/.default"}
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
