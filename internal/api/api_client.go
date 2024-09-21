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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers/array"
)

// ProviderClient is a wrapper around the API client that provides additional helper methods.
type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *Client
}

// GetConfig returns the provider configuration.
func (client *Client) GetConfig() *config.ProviderConfig {
	return client.Config
}

// Client is a base client for specific API clients implmented in services.
type Client struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}

// ApiHttpResponse is a wrapper around http.Response that provides additional helper methods.
func NewApiClientBase(providerConfig *config.ProviderConfig, baseAuth *Auth) *Client {
	return &Client{
		Config:   providerConfig,
		BaseAuth: baseAuth,
	}
}

var retryableStatusCodes = []int{
	http.StatusUnauthorized,    // 401 is retryable because the token may have expired.
	http.StatusRequestTimeout,  // 408 is retryable because the request may have timed out.
	http.StatusTooEarly,        // 425 is retryable because the request may have been rate limited.
	http.StatusTooManyRequests, // 429 is retryable because the request may have been rate limited.
	http.StatusInternalServerError,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
	http.StatusTooManyRequests,
}

// Execute executes an HTTP request with the given method, url, headers, and body.
func (client *Client) Execute(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*HttpResponse, error) {
	if len(scopes) == 0 {
		// if no scopes are provided, try to guess the scope from the URL.
		scope, err := tryGetScopeFromURL(url, client.Config.Urls)
		if err != nil {
			return nil, err
		}

		scopes = []string{scope}
	}

	if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
		return nil, helpers.WrapIntoProviderError(e, helpers.ERROR_INCORRECT_URL_FORMAT, "when using scope, the calling url must be an absolute url, not a relative path")
	}

	for {
		token, err := client.BaseAuth.GetTokenForScopes(ctx, scopes)
		if err != nil {
			return nil, err
		}

		bodyBuffer, err := prepareRequestBody(body)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
		if err != nil {
			return nil, err
		}

		resp, err := client.doRequest(ctx, token, request, headers)
		if err != nil {
			return resp, err
		}

		isAcceptable := len(acceptableStatusCodes) > 0 && array.Contains(acceptableStatusCodes, resp.Response.StatusCode)
		if isAcceptable {
			if responseObj != nil {
				err = resp.MarshallTo(responseObj)
				if err != nil {
					return resp, err
				}
			}

			return resp, nil
		}

		isRetryable := array.Contains(retryableStatusCodes, resp.Response.StatusCode)
		if !isRetryable {
			return resp, helpers.WrapIntoProviderError(err, helpers.ERROR_UNEXPECTED_HTTP_RETURN_CODE, fmt.Sprintf("expected status code: %d, recieved: [%d]", acceptableStatusCodes, resp.Response.StatusCode))
		}

		waitFor := retryAfter(ctx, resp.Response)

		tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", resp.Response.StatusCode, url, waitFor))

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return resp, err
		}
	}
}

// ExecuteForGivenScope executes an HTTP request with the given scope.
// The scope is used to obtain an access token for the request.
// The method, url, headers, and body are used to construct the request.
// The acceptableStatusCodes are used to validate the response status code.
// The responseObj is used to unmarshal the response body from json.
// If the responseObj is nil, the response body is not unmarshalled.
// If the response status code is not in the acceptableStatusCodes, an error is returned.
// func (client *Client) ExecuteForGivenScope(ctx context.Context, scope, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*HttpResponse, error) {
// 	if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
// 		return nil, helpers.WrapIntoProviderError(e, helpers.ERROR_INCORRECT_URL_FORMAT, "when using scope, the calling url must be an absolute url, not a relative path")
// 	}

// 	token, err := client.BaseAuth.GetTokenForScopes(ctx, []string{scope})
// 	if err != nil {
// 		return nil, err
// 	}

// 	bodyBuffer, err := prepareRequestBody(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
// 	if err != nil {
// 		return nil, err
// 	}

// 	apiResponse, err := client.doRequest(ctx, token, request, headers)
// 	if err != nil {
// 		return apiResponse, err
// 	}

// 	isStatusCodeValid := array.Contains(acceptableStatusCodes, apiResponse.Response.StatusCode)
// 	if !isStatusCodeValid && len(acceptableStatusCodes) > 0 {
// 		return apiResponse, helpers.WrapIntoProviderError(err, helpers.ERROR_UNEXPECTED_HTTP_RETURN_CODE, fmt.Sprintf("expected status code: %d, recieved: [%d]", acceptableStatusCodes, apiResponse.Response.StatusCode))
// 	}

// 	if responseObj != nil {
// 		err = apiResponse.MarshallTo(responseObj)
// 		if err != nil {
// 			return apiResponse, err
// 		}
// 	}
// 	return apiResponse, nil
// }

// RetryAfterDefault returns a random duration between 5 and 10 seconds.
func (client *Client) RetryAfterDefault() time.Duration {
	retryAfter5to10Seconds := time.Duration((rand.Intn(5) + 5)) * time.Second
	return retryAfter5to10Seconds
}

// SleepWithContext sleeps for the given duration or until the context is canceled.
func (client *Client) SleepWithContext(ctx context.Context, duration time.Duration) error {
	if helpers.IsTestContext(ctx) {
		return nil
	}
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

// TryGetScopeFromURL returns the authorization scope for the given API URL.
func tryGetScopeFromURL(url string, cloudConfig config.ProviderConfigUrls) (string, error) {
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

func prepareRequestBody(body any) (io.Reader, error) {
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

	return bodyBuffer, nil
}
