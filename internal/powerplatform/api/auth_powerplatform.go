package powerplatform

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (client *PowerPlatformAuth) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	//todo implement
	panic("[AuthenticateUserPass] not implemented")
}

func (client *PowerPlatformAuth) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	scopes := []string{"https://api.powerplatform.com/.default"}
	token, expiry, err := client.baseAuth.AuthClientSecret(ctx, scopes, tenantId, applicationId, secret)
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
