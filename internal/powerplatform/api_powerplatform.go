package powerplatform

import (
	"context"
	"errors"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

var _ PowerPlatformClientInterface = &PowerPlatformClientImplementation{}

type PowerPlatformClientInterface interface {
	Initialize(context.Context) (string, error)

	GetBillingPolicies(ctx context.Context) ([]models.BillingPoliciesDto, error)
}

type PowerPlatformClientImplementation struct {
	Config ProviderConfig
	Auth   PowerPlatformAuthInterface
}

func (client *PowerPlatformClientImplementation) Initialize(ctx context.Context) (string, error) {

	if client.Auth.IsTokenExpiredOrEmpty() {
		if client.Config.Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.Auth.AuthenticateClientSecret(ctx, client.Config.Credentials.TenantId, client.Config.Credentials.ClientId, client.Config.Credentials.Secret)
			if err != nil {
				return "", err
			}
			return token, nil
		} else if client.Config.Credentials.IsUserPassCredentialsProvided() {
			token, err := client.Auth.AuthenticateUserPass(ctx, client.Config.Credentials.TenantId, client.Config.Credentials.Username, client.Config.Credentials.Password)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			return "", errors.New("no credentials provided")
		}
	} else {
		//todo this is not implemented yet
		token, err := client.Auth.RefreshToken()
		if err != nil {
			return "", err
		}
		return token, nil

	}
}

func (client *PowerPlatformClientImplementation) GetBillingPolicies(ctx context.Context) ([]models.BillingPoliciesDto, error) {
	//todo implement
	return nil, nil
}
