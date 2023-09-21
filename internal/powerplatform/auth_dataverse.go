package powerplatform

import (
	"context"
	"strings"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

var _ DataverseAuthInterface = &DataverseAuthImplementation{}

type DataverseAuthInterface interface {
	IsTokenExpiredOrEmpty(environmentId string) bool
	RefreshToken(environmentId string) (string, error)

	AuthenticateUserPass(ctx context.Context, environmentId, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, environmentid, tenantId, applicationid, secret string) (string, error)
}

type DataverseAuthImplementation struct {
	tokensByEnvironmentId map[string]*DataverseAuth
}

type DataverseAuth struct {
	Token          string
	TokenExpiry    time.Time
	EnvironmentUrl string
}

func (client *DataverseAuthImplementation) setAuthDataInCache(environmentId string, authData *DataverseAuth) {
	if client.tokensByEnvironmentId == nil {
		client.tokensByEnvironmentId = make(map[string]*DataverseAuth)
	}
	client.tokensByEnvironmentId[environmentId] = authData
}

func (client *DataverseAuthImplementation) getAuthDataFromCache(environmentId string) *DataverseAuth {
	if client.tokensByEnvironmentId == nil {
		client.tokensByEnvironmentId = make(map[string]*DataverseAuth)
		return nil
	} else {
		auth, exist := client.tokensByEnvironmentId[environmentId]
		if exist {
			return auth
		} else {
			return nil
		}
	}
}

func (client *DataverseAuthImplementation) IsTokenExpiredOrEmpty(environmentId string) bool {
	auth := client.getAuthDataFromCache(environmentId)
	if auth == nil {
		return true
	} else {
		return time.Now().After(auth.TokenExpiry)
	}
}

func (client *DataverseAuthImplementation) RefreshToken(environmentId string) (string, error) {
	//todo implement token refresh
	panic("[RefreshToken] not implemented")
}

func (client *DataverseAuthImplementation) AuthenticateUserPass(ctx context.Context, environmentId, tenantId, username, password string) (string, error) {
	//todo implement when needed
	panic("[AuthenticateUserPass] not implemented")
}

func (client *DataverseAuthImplementation) AuthenticateClientSecret(ctx context.Context, environmentUrl, tenantId, applicationid, secret string) (string, error) {
	environmentUrl = strings.TrimSuffix(environmentUrl, "/")

	scopes := []string{environmentUrl + "//.default"}
	token, expiry, err := client.authClientSecret(ctx, scopes, tenantId, applicationid, secret)
	if err != nil {
		return "", err
	}

	client.setAuthDataInCache(environmentUrl, &DataverseAuth{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil
}

func (client *DataverseAuthImplementation) authClientSecret(ctx context.Context, scopes []string, tenantId, applicationId, clientSecret string) (string, time.Time, error) {
	authority := "https://login.microsoftonline.com/" + tenantId

	cred, err := confidential.NewCredFromSecret(clientSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	confidentialClientApp, err := confidential.New(authority, applicationId, cred)
	if err != nil {
		return "", time.Time{}, err
	}

	authResult, err := confidentialClientApp.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		return "", time.Time{}, err
	}

	return authResult.AccessToken, authResult.ExpiresOn, nil
}
