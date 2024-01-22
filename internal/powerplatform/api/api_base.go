package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
)

func (client *ApiClientBase) GetConfig() *config.ProviderConfig {
	return client.Config
}

type ApiClientBase struct {
	Config   *config.ProviderConfig
	BaseAuth *AuthBase
}

func NewApiClientBase(config *config.ProviderConfig, baseAuth *AuthBase) *ApiClientBase {
	return &ApiClientBase{
		Config:   config,
		BaseAuth: baseAuth,
	}
}

func (client *ApiClientBase) ExecuteBase(ctx context.Context, token, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	var bodyBuffer io.Reader = nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyBytes)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	apiResponse, err := client.doRequest(token, request, headers)
	if err != nil {
		return nil, err
	}

	isStatusCodeValid := false
	for _, statusCode := range acceptableStatusCodes {
		if apiResponse.Response.StatusCode == statusCode {
			isStatusCodeValid = true
			break
		}
	}
	if !isStatusCodeValid {
		return nil, fmt.Errorf("expected status code: %d, recieved: %d", acceptableStatusCodes, apiResponse.Response.StatusCode)
	}

	if responseObj != nil {
		err = apiResponse.MarshallTo(responseObj)
		if err != nil {
			return nil, err
		}
	}
	return apiResponse, nil
}

func (client *ApiClientBase) doRequest(token string, request *http.Request, headers http.Header) (*ApiHttpResponse, error) {
	apiHttpResponse := &ApiHttpResponse{}

	if headers != nil {
		request.Header = headers
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	//todo validate that initializing the http client everytime is ok from performance perspective
	httpClient := http.DefaultClient

	if request.Header["Authorization"] == nil {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	request.Header.Set("User-Agent", "terraform-provider-power-platform")

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = body
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(body))
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}
	return apiHttpResponse, nil
}

func (client *ApiClientBase) InitializeBase(ctx context.Context, scopes []string, auth AuthBaseOperationInterface) (string, error) {
	token, err := client.BaseAuth.GetToken()

	if _, ok := err.(*TokeExpiredError); ok {
		tflog.Debug(ctx, "Token expired. authenticating...")

		if client.Config.Credentials.IsClientSecretCredentialsProvided() {
			token, err := auth.AuthenticateClientSecret(ctx, scopes, client.Config.Credentials)
			if err != nil {
				return "", err
			}
			tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
			return token, nil
		} else if client.Config.Credentials.IsUserPassCredentialsProvided() {
			token, err := auth.AuthenticateUserPass(ctx, scopes, client.Config.Credentials)
			if err != nil {
				return "", err
			}
			tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
			return token, nil
		} else if client.Config.Credentials.UseCli {
			token, err := auth.AuthUsingCli(ctx, scopes, client.Config.Credentials)
			tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
			return token, err
		} else {
			return "", errors.New("no credentials provided")
		}

	} else if err != nil {
		return "", err
	} else {
		tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
		return token, nil
	}
}

func (client *ApiClientBase) InitializeBaseWithUserNamePassword(ctx context.Context, scopes []string, auth AuthBaseOperationInterface) (string, error) {
	token, err := client.BaseAuth.GetToken()

	if _, ok := err.(*TokeExpiredError); ok {
		tflog.Debug(ctx, "Token expired. authenticating...")

		token, err := auth.AuthenticateUserPass(ctx, scopes, client.Config.Credentials)
		if err != nil {
			return "", err
		}
		tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
		return token, nil

	} else if err != nil {
		return "", err
	} else {
		tflog.Debug(ctx, fmt.Sprintln("Token aquired: ", "********"))
		return token, nil
	}
}

type ApiHttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (apiResponse *ApiHttpResponse) MarshallTo(obj interface{}) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (apiResponse *ApiHttpResponse) GetHeader(name string) string {
	return apiResponse.Response.Header.Get(name)
}

func (ApiHttpResponse *ApiHttpResponse) ValidateStatusCode(expectedStatusCode int) error {
	if ApiHttpResponse.Response.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status code: %d, recieved: %d", expectedStatusCode, ApiHttpResponse.Response.StatusCode)
	}
	return nil
}
