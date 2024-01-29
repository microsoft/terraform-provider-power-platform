package powerplatform

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
)

var _ AuthBaseOperationInterface = &DataverseAuth{}

type DataverseAuth struct {
	baseAuth               *AuthBase
	tokensByEnvironmentUrl map[string]*DataverseAuthDetails
}

func NewDataverseAuth(baseAuth *AuthBase) *DataverseAuth {
	return &DataverseAuth{
		baseAuth: baseAuth,
	}
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

func (client *DataverseAuth) GetToken(scopes []string) (string, error) {
	isExpired, token := client.isTokenExpiredOrEmpty(strings.Join(scopes, "_"))
	if isExpired {
		return "", &TokeExpiredError{Message: "token is expired or empty"}
	} else {
		return token, nil
	}
}

func (client *DataverseAuth) isTokenExpiredOrEmpty(environmentUrl string) (bool, string) {
	auth := client.getAuthDataFromCache(environmentUrl)
	isExpired := auth == nil || (auth != nil && auth.Token == "") || (auth != nil && time.Now().After(auth.TokenExpiry))
	if isExpired {
		return true, ""
	} else {
		return false, auth.Token
	}
}

func (client *DataverseAuth) AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthUsingCli(ctx, scopes, credentials)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your client/secret and try again")
		}
		return "", err
	}
	client.setAuthDataInCache(strings.Join(scopes, "_"), &DataverseAuthDetails{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil

}

func (client *DataverseAuth) AuthenticateUserPass(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthenticateUserPass(ctx, scopes, credentials)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.setAuthDataInCache(strings.Join(scopes, "_"), &DataverseAuthDetails{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil
}

func (client *DataverseAuth) AuthenticateClientSecret(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthClientSecret(ctx, scopes, credentials)
	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.setAuthDataInCache(strings.Join(scopes, "_"), &DataverseAuthDetails{
		Token:       token,
		TokenExpiry: expiry,
	})
	return token, nil
}
