// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
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

func (client *ApiClient) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	scope, err := TryGetScopeFromURL(url, client.Config.Urls)

	if err != nil {
		return nil, err
	}

	token, err := client.BaseAuth.GetTokenForScopes(ctx, []string{scope})
	if err != nil {
		return nil, err
	}

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
	if len(acceptableStatusCodes) == 0 {
		isStatusCodeValid = true
	} else {
		for _, statusCode := range acceptableStatusCodes {
			if apiResponse.Response.StatusCode == statusCode {
				isStatusCodeValid = true
				break
			}
		}
		if !isStatusCodeValid {
			return nil, fmt.Errorf("expected status code: %d, recieved: %d", acceptableStatusCodes, apiResponse.Response.StatusCode)
		}
	}

	if responseObj != nil {
		err = apiResponse.MarshallTo(responseObj)
		if err != nil {
			return nil, err
		}
	}
	return apiResponse, nil
}

func (client *ApiClient) Sleep(duration time.Duration) {
	if client.Config.Credentials.TestMode {
		//Don't sleep during testing
		return
	} else {
		//lintignore:R018
		time.Sleep(duration)
	}
}
