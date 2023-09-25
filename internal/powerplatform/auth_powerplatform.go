package powerplatform

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

var _ PowerPlatformAuthInterface = &PowerPlatformAuthImplementation{}

type PowerPlatformAuthInterface interface {
	GetBase() common.AuthInterface

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type PowerPlatformAuthImplementation struct {
	BaseAuth common.AuthInterface
}

func (client *PowerPlatformAuthImplementation) GetBase() common.AuthInterface {
	return client.BaseAuth
}

func (client *PowerPlatformAuthImplementation) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	//todo implement
	panic("[AuthenticateUserPass] not implemented")
}

func (client *PowerPlatformAuthImplementation) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
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
