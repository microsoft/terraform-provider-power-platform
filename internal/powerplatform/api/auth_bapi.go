package powerplatform

import (
	"context"

	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
)

var _ AuthBaseOperationInterface = &BapiAuth{}

type BapiAuth struct {
	baseAuth *AuthBase
}

func NewBapiAuth(authBase *AuthBase) *BapiAuth {
	return &BapiAuth{
		baseAuth: authBase,
	}
}

func (client *BapiAuth) GetBase() *AuthBase {
	return client.baseAuth
}

func (client *BapiAuth) AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	panic("to remove")
	// token, expiry, err := client.baseAuth.AuthUsingCli(ctx, scopes, credentials)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
	// 		tflog.Debug(ctx, err.Error())
	// 		return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
	// 	}
	// 	return "", err
	// }
	// client.baseAuth.SetToken(token)
	// client.baseAuth.SetTokenExpiry(expiry)

	// return client.baseAuth.GetToken()
}

func (client *BapiAuth) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	panic("to remove")
	// token, expiry, err := client.baseAuth.AuthenticateUserPass(ctx, scopes, credentials)

	// if err != nil {
	// 	if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
	// 		tflog.Debug(ctx, err.Error())
	// 		return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
	// 	}
	// 	return "", err
	// }

	// client.baseAuth.SetToken(token)
	// client.baseAuth.SetTokenExpiry(expiry)

	// return client.baseAuth.GetToken()
}

func (client *BapiAuth) AuthenticateClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	panic("to remove")
	// token, expiry, err := client.baseAuth.AuthClientSecret(ctx, scopes, credentials)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
	// 		tflog.Debug(ctx, err.Error())
	// 		return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
	// 	}
	// 	return "", err
	// }

	// client.baseAuth.SetToken(token)
	// client.baseAuth.SetTokenExpiry(expiry)

	// return client.baseAuth.GetToken()
}
