package powerplatform

import (
	"context"
	"strings"
	"time"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ DataverseAuthInterface = &DataverseAuth{}

type DataverseAuthInterface interface {
	GetToken(environmentUrl string) (string, error)
	AuthenticateUserPass(ctx context.Context, environmentUrl, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, environmentUrl, tenantId, applicationid, secret string) (string, error)
}

type DataverseAuth struct {
	BaseAuth               api.AuthInterface
	tokensByEnvironmentUrl map[string]*DataverseAuthDetails
}

type DataverseAuthDetails struct {
	Token          string
	TokenExpiry    time.Time
	EnvironmentUrl string
}

func (client *DataverseAuth) setAuthDataInCache(environmentUrl string, authData *DataverseAuthDetails) {
	if client.tokensByEnvironmentUrl == nil {
		client.tokensByEnvironmentUrl = make(map[string]*DataverseAuthDetails)
	}
	client.tokensByEnvironmentUrl[environmentUrl] = authData
}

func (client *DataverseAuth) getAuthDataFromCache(environmentUrl string) *DataverseAuthDetails {
	if client.tokensByEnvironmentUrl == nil {
		client.tokensByEnvironmentUrl = make(map[string]*DataverseAuthDetails)
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

func (client *DataverseAuth) GetToken(environmentUrl string) (string, error) {
	if client.isTokenExpiredOrEmpty(environmentUrl) {
		return "", &api.TokeExpiredError{Message: "token is expired or empty"}
	}
	return client.getAuthDataFromCache(environmentUrl).Token, nil
}

func (client *DataverseAuth) isTokenExpiredOrEmpty(environmentUrl string) bool {
	auth := client.getAuthDataFromCache(environmentUrl)
	return auth == nil || (auth != nil && auth.Token == "") || (auth != nil && time.Now().After(auth.TokenExpiry))
}

func (client *DataverseAuth) AuthenticateUserPass(ctx context.Context, environmentUrl, tenantId, username, password string) (string, error) {
	//todo implement when needed
	panic("[AuthenticateUserPass] not implemented")
}

func (client *DataverseAuth) AuthenticateClientSecret(ctx context.Context, environmentUrl, tenantId, applicationid, secret string) (string, error) {
	environmentUrl = strings.TrimSuffix(environmentUrl, "/")

	scopes := []string{environmentUrl + "//.default"}
	token, expiry, err := client.BaseAuth.AuthClientSecret(ctx, scopes, tenantId, applicationid, secret)
	if err != nil {
		return "", err
	}

	client.setAuthDataInCache(environmentUrl, &DataverseAuthDetails{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil
}
