package powerplatform

import (
	"context"
	"time"

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
	config *config.ProviderConfig
	//token       string
	//tokenExpiry time.Time
	memCache  *common.MemoryCache
	authCache *common.AuthenticationCache
}

func NewAuthBase(config *config.ProviderConfig) *AuthBase {
	return &AuthBase{
		config:    config,
		authCache: common.NewAuthenticationCache(),
		memCache:  common.NewMemoryCache(),
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
	//TODO if use used different clientid in cli this will not work!?
	publicClient, err := public.New(constants.CLIENT_ID, public.WithCache(client.authCache))
	if err != nil {
		return "", time.Time{}, err
	}

	accounts, err := publicClient.Accounts(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	if len(accounts) == 0 {
		return "", time.Time{}, &TokeExpiredError{"No cached tokens found. Please login using 'CLI' first and try again."}
	}

	//TODO cache file may have many accounts maybe we shoudn't use the first one but the one that matches the tenantid or username from credentials given in provider by the user?
	credentials.TenantId = accounts[0].Realm
	authResult, err := publicClient.AcquireTokenSilent(ctx, scopes, public.WithTenantID(credentials.TenantId), public.WithSilentAccount(accounts[0]))
	if err != nil {
		return "", time.Time{}, err
	}
	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {

	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority(client.GetAuthority(credentials.TenantId)), public.WithCache(client.memCache))
	if err != nil {
		return "", time.Time{}, err
	}

	accounts, err := publicClient.Accounts(ctx)
	if err != nil {
		return "", time.Time{}, err
	}

	savedAccount := public.Account{}
	for _, account := range accounts {
		if account.PreferredUsername == credentials.Username {
			savedAccount = account
			break
		}
	}

	if savedAccount.PreferredUsername != "" {
		authResult, err := publicClient.AcquireTokenSilent(ctx, scopes, public.WithTenantID(credentials.TenantId), public.WithSilentAccount(savedAccount))
		if err != nil {
			return "", time.Time{}, err
		}
		return authResult.AccessToken, authResult.ExpiresOn, nil
	} else {
		authResult, err := publicClient.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)
		if err != nil {
			return "", time.Time{}, err
		}
		return authResult.AccessToken, authResult.ExpiresOn, nil
	}

	// publicClientApp, err := public.New(constants.CLIENT_ID, public.WithAuthority(client.GetAuthority(credentials.TenantId)))
	// if err != nil {
	// 	return "", time.Time{}, err
	// }

	// authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)
	// if err != nil {
	// 	return "", time.Time{}, err
	// }

	// return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *AuthBase) AuthClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {
	panic("to remove")

	// cred, err := confidential.NewCredFromSecret(credentials.Secret)
	// if err != nil {
	// 	return "", time.Time{}, err
	// }

	// confidentialClientApp, err := confidential.New(client.GetAuthority(credentials.TenantId), credentials.ClientId, cred)
	// if err != nil {
	// 	return "", time.Time{}, err
	// }

	// authResult, err := confidentialClientApp.AcquireTokenByCredential(ctx, scopes)
	// if err != nil {
	// 	return "", time.Time{}, err
	// }

	// return authResult.AccessToken, authResult.ExpiresOn, nil
}

// func (client *AuthBase) SetToken(token string) {
// 	client.token = token
// }

// func (client *AuthBase) SetTokenExpiry(tokenExpiry time.Time) {
// 	client.tokenExpiry = tokenExpiry
// }

// func (client *AuthBase) GetTokenExpiry() time.Time {
// 	return client.tokenExpiry
// }

// func (client *AuthBase) GetToken() (string, error) {
// 	if client.IsTokenExpiredOrEmpty() {
// 		return "", &TokeExpiredError{"token is expired or empty"}
// 	}
// 	return client.token, nil
// }

// func (client *AuthBase) IsTokenExpiredOrEmpty() bool {
// 	return client.token == "" || time.Now().After(client.tokenExpiry)
// }
