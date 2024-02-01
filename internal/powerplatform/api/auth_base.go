package powerplatform

import (
	"context"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	common "github.com/microsoft/terraform-provider-power-platform/common"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
)

type TokeExpiredError struct {
	Message string
}

func (e *TokeExpiredError) Error() string {
	return e.Message
}

type AuthBase struct {
	config      *config.ProviderConfig
	token       string
	tokenExpiry time.Time
	authCache   *common.AuthenticationCache
}

func NewAuthBase(config *config.ProviderConfig) *AuthBase {
	return &AuthBase{
		config:    config,
		authCache: common.NewAuthenticationCache(),
	}
}

type AuthBaseOperationInterface interface {
	AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error)
	AuthenticateClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error)
	AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error)
}

func (client *AuthBase) GetAuthority(tenantid string) string {
	return "https://login.microsoftonline.com/" + tenantid
}

func (client *AuthBase) AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithCache(client.authCache))
	if err != nil {
		return "", time.Time{}, err
	}

	defaultAccount, err := client.authCache.GetDefaultAccount(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	credentials.TenantId = defaultAccount.Realm
	authResult, err := publicClient.AcquireTokenSilent(ctx, scopes, public.WithTenantID(credentials.TenantId), public.WithSilentAccount(*defaultAccount))
	if err != nil {
		return "", time.Time{}, err
	}
	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {
	publicClientApp, err := public.New(constants.CLIENT_ID, public.WithAuthority(client.GetAuthority(credentials.TenantId)))
	if err != nil {
		return "", time.Time{}, err
	}

	authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)
	if err != nil {
		return "", time.Time{}, err
	}

	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) AuthClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {

	cred, err := confidential.NewCredFromSecret(credentials.Secret)
	if err != nil {
		return "", time.Time{}, err
	}

	confidentialClientApp, err := confidential.New(client.GetAuthority(credentials.TenantId), credentials.ClientId, cred)
	if err != nil {
		return "", time.Time{}, err
	}

	authResult, err := confidentialClientApp.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		return "", time.Time{}, err
	}

	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) SetToken(token string) {
	client.token = token
}

func (client *AuthBase) SetTokenExpiry(tokenExpiry time.Time) {
	client.tokenExpiry = tokenExpiry
}

func (client *AuthBase) GetTokenExpiry() time.Time {
	return client.tokenExpiry
}

func (client *AuthBase) GetToken() (string, error) {
	if client.IsTokenExpiredOrEmpty() {
		return "", &TokeExpiredError{"token is expired or empty"}
	}
	return client.token, nil
}

func (client *AuthBase) IsTokenExpiredOrEmpty() bool {
	return client.token == "" || time.Now().After(client.tokenExpiry)
}
