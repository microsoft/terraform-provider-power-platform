package powerplatform_bapi

import (
	"context"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type AuthResponse struct {
	Token string `json:"token"`
}
type AuthType int64

const (
	UsernamePassword AuthType = 0
	ServicePrincipal AuthType = 1
)

func (client *ApiClient) DoAuthUsernamePassword(ctx context.Context, tenantId, username, password string) (*AuthResponse, error) {

	scopes := []string{"https://service.powerapps.com//.default"}
	publicClientApplicationID := "1950a258-227b-4e31-a9cf-717495945fc2"
	authority := "https://login.microsoftonline.com/" + tenantId

	publicClientApp, err := public.New(publicClientApplicationID, public.WithAuthority(authority))
	if err != nil {
		return nil, err
	}

	authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, username, password)
	if err != nil {
		return nil, err
	}

	client.Provider.TenantId = tenantId
	client.Provider.Username = username
	client.Provider.Password = password
	client.Token = authResult.AccessToken

	authResponse := AuthResponse{
		Token: authResult.AccessToken,
	}
	return &authResponse, nil
}

func (client *ApiClient) DoAuthClientSecret(ctx context.Context, tenantId, applicationId, clientSecret string) (*AuthResponse, error) {
	scopes := []string{"https://service.powerapps.com//.default"}
	client.Provider.TenantId = tenantId
	client.Provider.ClientId = applicationId
	client.Provider.ClientSecret = clientSecret
	auth, err := client.authClientSecret(ctx, scopes, tenantId, applicationId, clientSecret)
	if err != nil {
		return nil, err
	}
	client.Token = auth.Token
	return auth, nil
}

func (client *ApiClient) DoAuthClientSecretForDataverse(ctx context.Context, environmentUrl string) (*AuthResponse, error) {

	environmentUrl = strings.TrimSuffix(environmentUrl, "/")

	//todo look at token's expiration time and refresh if needed
	//we don't expect terraform operation to run long enough for it to expire though
	if auth, ok := client.DataverseAuthMap[environmentUrl]; ok {
		return auth, nil
	} else {
		scopes := []string{environmentUrl + "//.default"}
		auth, err := client.authClientSecret(ctx, scopes, client.Provider.TenantId, client.Provider.ClientId, client.Provider.ClientSecret)
		if err != nil {
			return nil, err
		}
		client.DataverseAuthMap[environmentUrl] = auth
		return auth, nil
	}
}

func (client *ApiClient) authClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (*AuthResponse, error) {
	authority := "https://login.microsoftonline.com/" + tenantId

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		return nil, err
	}

	confidentialClientApp, err := confidential.New(authority, applicationId, cred)
	if err != nil {
		return nil, err
	}

	authResult, err := confidentialClientApp.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		return nil, err
	}

	authResponse := AuthResponse{
		Token: authResult.AccessToken,
	}
	return &authResponse, err
}
