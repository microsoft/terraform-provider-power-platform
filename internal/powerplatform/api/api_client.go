// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"reflect"
	"strings"
	"time"

	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *ApiClient
}

func (client *ApiClient) GetConfig() *config.ProviderConfig {
	return client.Config
}

type ApiClient struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}

func NewApiClientBase(config *config.ProviderConfig, baseAuth *Auth) *ApiClient {
	return &ApiClient{
		Config:   config,
		BaseAuth: baseAuth,
	}
}

func TryGetScopeFromURL(url string, cloudConfig config.ProviderConfigUrls) (string, error) {

	switch {
	case strings.LastIndex(url, cloudConfig.BapiUrl) != -1,
		strings.LastIndex(url, cloudConfig.PowerAppsUrl) != -1:
		return cloudConfig.PowerAppsScope, nil

	case strings.LastIndex(url, cloudConfig.PowerPlatformUrl) != -1:
		return cloudConfig.PowerPlatformScope, nil

	default:
		u, err := neturl.Parse(url)
		return u.Scheme + "://" + u.Host + "/.default", err
	}
}

func (client *ApiClient) ExecuteForGivenScope(ctx context.Context, scope, method, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, helpers.WrapIntoProviderError(nil, helpers.ERROR_INCORRECT_URL_FORMAT, "when using scope, the calling url must be an absolute url, not a relative path")
	}
	token, err := client.BaseAuth.GetTokenForScopes(ctx, []string{scope})
	if err != nil {
		return nil, err
	}

	var bodyBuffer io.Reader = nil
	if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
		if reflect.ValueOf(body).Kind() == reflect.Ptr && reflect.ValueOf(body).Elem().Kind() == reflect.String {
			strp, _ := body.(*string)
			bodyBuffer = strings.NewReader(*strp)
		} else {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyBuffer = bytes.NewBuffer(bodyBytes)
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	apiResponse, err := client.doRequest(token, request, headers)
	if err != nil {
		return apiResponse, err
	}

	isStatusCodeValid := false
	if len(acceptableStatusCodes) == 0 {
		isStatusCodeValid = true
	} else {
		for _, statusCode := range acceptableStatusCodes {
			if apiResponse.Response.StatusCode == statusCode {
				isStatusCodeValid = true
				break
			}
		}
	}

	if !isStatusCodeValid {
		return apiResponse, helpers.WrapIntoProviderError(err, helpers.ERROR_UNEXPECTED_HTTP_RETURN_CODE, fmt.Sprintf("expected status code: %d, recieved: [%d]", acceptableStatusCodes, apiResponse.Response.StatusCode))
	}

	if responseObj != nil {
		err = apiResponse.MarshallTo(responseObj)
		if err != nil {
			return apiResponse, err
		}
	}
	return apiResponse, nil
}

func (client *ApiClient) Execute(ctx context.Context, method, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	scope, err := TryGetScopeFromURL(url, client.Config.Urls)
	if err != nil {
		return nil, err
	}
	return client.ExecuteForGivenScope(ctx, scope, method, url, headers, body, acceptableStatusCodes, responseObj)
}

func (client *ApiClient) SleepWithContext(ctx context.Context, duration time.Duration) error {
	if client.Config.Credentials.TestMode {
		//Don't sleep during testing
		return nil
	} else {
		select {
		case <-time.After(duration):
			// Time has elapsed
			return nil
		case <-ctx.Done():
			// Context was canceled
			return ctx.Err()
		}
	}
}
