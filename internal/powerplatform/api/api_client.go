// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func TryGetScopeFromURL(url string) (string, error) {
	switch {
	case strings.LastIndex(url, "api.bap.microsoft.com") != -1,
		strings.LastIndex(url, "api.powerapps.com") != -1:

		return "https://service.powerapps.com/.default", nil
	case strings.LastIndex(url, "api.powerplatform.com") != -1:

		return "https://api.powerplatform.com/.default", nil
	case strings.LastIndex(url, ".com/") != -1:

		scope := strings.SplitAfterN(url, ".com/", 2)[0]
		scope = scope + ".default"
		return scope, nil
	default:
		return "", errors.New("Unable to determine scope from url: '" + url + "'. Please provide your own scope.")
	}
}

func (client *ApiClient) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	scope, err := TryGetScopeFromURL(url)
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