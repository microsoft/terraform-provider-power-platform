package powerplatform

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

type Auth struct {
	config        *config.ProviderConfig
	fileCache     *common.FileCache
	memoryCache   *common.MemoryCache
	homeAccountID string
}

func NewAuthBase(config *config.ProviderConfig) *Auth {
	return &Auth{
		config:      config,
		fileCache:   common.NewAuthenticationCache(),
		memoryCache: common.NewMemoryCache(),
	}
}

func (client *Auth) GetAuthority(tenantid string) string {
	return constants.OAUTH_AUTHORITY_URL + tenantid
}

func (client *Auth) AuthenticateUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithCache(client.fileCache))
	if err != nil {
		return "", time.Time{}, err
	}

	defaultAccount, err := client.fileCache.GetDefaultAccount(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	if defaultAccount == nil {
		return "", time.Time{}, errors.New("no default account found. Please login CLI using 'terraform-provider-power-platform login' command")
	}

	credentials.TenantId = defaultAccount.Realm
	authResult, err := publicClient.AcquireTokenSilent(ctx, scopes, public.WithTenantID(credentials.TenantId), public.WithSilentAccount(*defaultAccount))
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", time.Time{}, errors.New("there was an issue authenticating with the provided credentials. Please check the your credentials and try again")
		}
		return "", time.Time{}, err
	}
	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *Auth) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority(client.GetAuthority(credentials.TenantId)), public.WithCache(client.memoryCache))
	if err != nil {
		return "", time.Time{}, err
	}

	authResult := public.AuthResult{}
	accounts, err := client.memoryCache.GetAccounts(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	if len(accounts) > 0 {
		authResult, err = publicClient.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount((accounts[len(accounts)-1])))

	} else {
		authResult, err = publicClient.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)
	}

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", time.Time{}, errors.New("there was an issue authenticating with the provided credentials. Please check the your credentials and try again")
		}
		return "", time.Time{}, err
	}
	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *Auth) AuthenticateClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, time.Time, error) {

	cred, err := confidential.NewCredFromSecret(credentials.Secret)
	if err != nil {
		return "", time.Time{}, err
	}
	confidentialClient, err := confidential.New(client.GetAuthority(credentials.TenantId), credentials.ClientId, cred, confidential.WithCache(client.memoryCache))
	if err != nil {
		return "", time.Time{}, err
	}

	authResult := confidential.AuthResult{}
	account, err := confidentialClient.Account(ctx, client.homeAccountID)
	if err != nil {
		return "", time.Time{}, err
	}

	if account.IsZero() {
		authResult, err = confidentialClient.AcquireTokenByCredential(ctx, scopes)
	} else {
		authResult, err = confidentialClient.AcquireTokenSilent(ctx, scopes, confidential.WithSilentAccount(account))
	}
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", time.Time{}, errors.New("there was an issue authenticating with the provided credentials. Please check the your credentials and try again")
		}
		return "", time.Time{}, err
	}
	//todo this doesn't work always correctly
	client.homeAccountID = fmt.Sprintf("-login.microsoftonline.com-accesstoken-%s-%s-%s", credentials.ClientId, credentials.TenantId, authResult.GrantedScopes[0])
	return authResult.AccessToken, authResult.ExpiresOn, nil
}

func (client *Auth) InitializeRequiredScopes(ctx context.Context, scopes []string) (string, error) {
	token := ""
	tokenExpiry := time.Time{}
	var err error

	switch {
	case client.config.Credentials.IsClientSecretCredentialsProvided():
		//todo use local credentials instead dependency injection
		token, tokenExpiry, err = client.AuthenticateClientSecret(ctx, scopes, client.config.Credentials)
	case client.config.Credentials.IsUserPassCredentialsProvided():
		token, tokenExpiry, err = client.AuthenticateUserPass(ctx, scopes, client.config.Credentials)
	case client.config.Credentials.IsCliProvided():
		token, tokenExpiry, err = client.AuthenticateUsingCli(ctx, scopes, client.config.Credentials)
	default:
		return "", errors.New("no credentials provided")
	}

	if err != nil {
		return "", err
	}
	tflog.Debug(ctx, fmt.Sprintf("Token acquired (expire: %s): **********", tokenExpiry))
	return token, nil
}

func (client *Auth) GetTokenForScope(ctx context.Context, scope string) (*string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[GetTokenForScope] Getting token for scope: '%s'", scope))

	token, err := client.InitializeRequiredScopes(ctx, []string{scope})
	return &token, err
}
