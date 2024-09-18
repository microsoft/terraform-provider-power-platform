// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	neturl "net/url"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// ProviderClient is a wrapper around the API client that provides additional helper methods.
type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *Client
}

// GetConfig returns the provider configuration.
func (client *ApiClient) GetConfig() *config.ProviderConfig {
	return client.Config
}

// ApiClient is a base client for specific API clients implmented in services.
type ApiClient struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}

// ApiHttpResponse is a wrapper around http.Response that provides additional helper methods.
func NewApiClientBase(config *config.ProviderConfig, baseAuth *Auth) *ApiClient {
	return &ApiClient{
		Config:   config,
		BaseAuth: baseAuth,
	}
}

// TryGetScopeFromURL returns the authorization scope for the given API URL.
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

// ExecuteForGivenScope executes an HTTP request with the given scope.
// The scope is used to obtain an access token for the request.
// The method, url, headers, and body are used to construct the request.
// The acceptableStatusCodes are used to validate the response status code.
// The responseObj is used to unmarshal the response body from json.
// If the responseObj is nil, the response body is not unmarshalled.
// If the response status code is not in the acceptableStatusCodes, an error is returned.
func (client *ApiClient) ExecuteForGivenScope(ctx context.Context, scope, method, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
		return nil, helpers.WrapIntoProviderError(e, helpers.ERROR_INCORRECT_URL_FORMAT, "when using scope, the calling url must be an absolute url, not a relative path")
	}

	token, err := client.BaseAuth.GetTokenForScopes(ctx, []string{scope})
	if err != nil {
		return nil, err
	}

	var bodyBuffer io.Reader
	if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
		if strp, ok := body.(*string); ok {
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
	apiResponse, err := client.doRequest(ctx, token, request, headers)
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

// Execute executes an HTTP request with the given method, url, headers, and body.
func (client *ApiClient) Execute(ctx context.Context, method, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	scope, err := TryGetScopeFromURL(url, client.Config.Urls)
	if err != nil {
		return nil, err
	}

	var response *HttpResponse
	for {
		response, err = client.ExecuteForGivenScope(ctx, scope, method, url, headers, body, acceptableStatusCodes, responseObj)
		if response == nil || response.Response == nil {
			return response, err
		}

		if response.Response.StatusCode != http.StatusUnauthorized &&
			response.Response.StatusCode != http.StatusGatewayTimeout &&
			response.Response.StatusCode != http.StatusTooManyRequests &&
			response.Response.StatusCode != http.StatusServiceUnavailable &&
			response.Response.StatusCode != http.StatusBadGateway {
			return response, err
		}

		defaultRetry := client.RetryAfterDefault()
		tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", response.Response.StatusCode, url, defaultRetry))

		err = client.SleepWithContext(ctx, defaultRetry)
		if err != nil {
			return response, err
		}
	}
}

// RetryAfterDefault returns a random duration between 5 and 10 seconds.
func (client *ApiClient) RetryAfterDefault() time.Duration {
	retryAfter5to10Seconds := time.Duration((rand.Intn(5) + 5)) * time.Second
	return retryAfter5to10Seconds
}

// SleepWithContext sleeps for the given duration or until the context is canceled.
func (client *ApiClient) SleepWithContext(ctx context.Context, duration time.Duration) error {
	if client.Config.TestMode {
		// Don't sleep during testing.
		return nil
	}
	select {
	case <-time.After(duration):
		// Time has elapsed.
		return nil
	case <-ctx.Done():
		// Context was canceled.
		return ctx.Err()
	}
}
