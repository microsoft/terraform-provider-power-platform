package powerplatform

import (
	"context"
	"errors"
	"strings"
	"time"

<<<<<<< HEAD
	"github.com/hashicorp/terraform-plugin-log/tflog"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
=======
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
>>>>>>> origin/main
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

<<<<<<< HEAD
func (client *DataverseAuth) AuthUsingCli(ctx context.Context, scopes []string, credentials *config.ProviderCredentials) (string, error) {
	token, expiry, err := client.baseAuth.AuthUsingCli(ctx, scopes, credentials)
=======
func (client *DataverseAuth) AuthenticateUserPass(ctx context.Context, environmentUrl string, credentials *powerplatform_common.ProviderCredentials) (string, error) {
	environmentUrl = strings.TrimSuffix(environmentUrl, "/")

	scopes := []string{environmentUrl + "//.default"}
	publicClientApplicationID := "1950a258-227b-4e31-a9cf-717495945fc2"
	authority := "https://login.microsoftonline.com/" + credentials.TenantId

	publicClientApp, err := public.New(publicClientApplicationID, public.WithAuthority(authority))

	if err != nil {
		return "", err
	}

	authResult, err := publicClientApp.AcquireTokenByUsernamePassword(ctx, scopes, credentials.Username, credentials.Password)

	if err != nil {
		if strings.Contains(err.Error(), "unable to resolve an endpoint: json decode error") {
			tflog.Debug(ctx, err.Error())
			return "", errors.New("there was an issue authenticating with the provided credentials. Please check the your username/password and try again")
		}
		return "", err
	}

	client.setAuthDataInCache(environmentUrl, &DataverseAuthDetails{
		Token:       authResult.AccessToken,
		TokenExpiry: authResult.ExpiresOn,
	})
	return authResult.AccessToken, nil
}
>>>>>>> origin/main

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
