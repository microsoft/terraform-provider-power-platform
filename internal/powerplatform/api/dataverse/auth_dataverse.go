package powerplatform

import (
	"context"
	"strings"
	"time"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ DataverseAuthInterface = &DataverseAuthImplementation{}

type DataverseAuthInterface interface {
	GetToken(environmentUrl string) (string, error)
	AuthenticateUserPass(ctx context.Context, environmentUrl, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, environmentUrl, tenantId, applicationid, secret string) (string, error)
}

type DataverseAuthImplementation struct {
	BaseAuth               api.AuthInterface
	tokensByEnvironmentUrl map[string]*DataverseAuth
}

type DataverseAuth struct {
	Token          string
	TokenExpiry    time.Time
	EnvironmentUrl string
}

func (client *DataverseAuthImplementation) setAuthDataInCache(environmentUrl string, authData *DataverseAuth) {
	if client.tokensByEnvironmentUrl == nil {
		client.tokensByEnvironmentUrl = make(map[string]*DataverseAuth)
	}
	client.tokensByEnvironmentUrl[environmentUrl] = authData
}

func (client *DataverseAuthImplementation) getAuthDataFromCache(environmentUrl string) *DataverseAuth {
	if client.tokensByEnvironmentUrl == nil {
		client.tokensByEnvironmentUrl = make(map[string]*DataverseAuth)
		return nil
	} else {
		auth, exist := client.tokensByEnvironmentUrl[environmentUrl]
		if exist {
			return auth
		} else {
			return nil
		}
	}
}

func (client *DataverseAuthImplementation) GetToken(environmentUrl string) (string, error) {
	if client.isTokenExpiredOrEmpty(environmentUrl) {
		return "", &api.TokeExpiredError{Message: "token is expired or empty"}
	}
	return client.getAuthDataFromCache(environmentUrl).Token, nil
}

func (client *DataverseAuthImplementation) isTokenExpiredOrEmpty(environmentUrl string) bool {
	auth := client.getAuthDataFromCache(environmentUrl)
	return auth == nil || (auth != nil && auth.Token == "") || (auth != nil && time.Now().After(auth.TokenExpiry))
}

func (client *DataverseAuthImplementation) AuthenticateUserPass(ctx context.Context, environmentUrl, tenantId, username, password string) (string, error) {
	//todo implement when needed
	panic("[AuthenticateUserPass] not implemented")
}

func (client *DataverseAuthImplementation) AuthenticateClientSecret(ctx context.Context, environmentUrl, tenantId, applicationid, secret string) (string, error) {
	environmentUrl = strings.TrimSuffix(environmentUrl, "/")

	scopes := []string{environmentUrl + "//.default"}
	token, expiry, err := client.BaseAuth.AuthClientSecret(ctx, scopes, tenantId, applicationid, secret)
	if err != nil {
		return "", err
	}

	client.setAuthDataInCache(environmentUrl, &DataverseAuth{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil
}
