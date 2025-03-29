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
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
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
	http.StatusUnauthorized,        // 401 is retryable because the token may have expired.
	http.StatusRequestTimeout,      // 408 is retryable because the request may have timed out.
	http.StatusTooEarly,            // 425 is retryable because the request may have been rate limited.
	http.StatusTooManyRequests,     // 429 is retryable because the request may have been rate limited.
	http.StatusInternalServerError, // 500 is retryable because the server may be overloaded.
	http.StatusBadGateway,          // 502 is retryable because the server may be overloaded.
	http.StatusServiceUnavailable,  // 503 is retryable because the server may be overloaded.
	http.StatusGatewayTimeout,      // 504 is retryable because the server may be overloaded.
}

// CaePolicyViolationError represents an error when a CAE policy violation is detected.
type CaePolicyViolationError struct {
	Message    string
	StatusCode int
	Headers    http.Header
}

func (e *CaePolicyViolationError) Error() string {
	return fmt.Sprintf("Continuous Access Evaluation policy violation: %s", e.Message)
}

// IsCaeChallengeResponse detects if a response contains a CAE challenge.
func IsCaeChallengeResponse(resp *http.Response) bool {
	if resp == nil {
		return false
	}

	// Check response headers for CAE challenge indicators.
	if resp.StatusCode == http.StatusUnauthorized {
		wwwAuthenticate := resp.Header.Get("WWW-Authenticate")
		if wwwAuthenticate != "" {
			return strings.Contains(wwwAuthenticate, "claims") &&
				strings.Contains(wwwAuthenticate, "access_denied")
		}
	}

	return false
}

// Execute executes an HTTP request with the given method, url, headers, and body.
//
// Parameters:
//   - ctx: context.Context - Provides context for the request, allowing for timeout and cancellation control.
//   - scopes: []string - A list of scopes that the request should be associated with. If no scopes are provided, the method attempts to infer the scope from the URL.
//   - method: string - Specifies the HTTP method to be used for the request (e.g., "GET", "POST", "PATCH").
//   - url: string - The URL to which the request is sent. This includes the scheme, host, path, and query parameters. The URL must be absolute and properly formatted.
//   - headers: http.Header - A collection of HTTP headers to include in the request. Headers provide additional information about the request, such as content type, authorization tokens, and custom metadata.
//   - body: any - The body of the request, which can be of any type. This is typically used for methods like POST and PATCH, where data needs to be sent to the server.
//   - acceptableStatusCodes: []int - A list of HTTP status codes that are considered acceptable for the response. If the response status code is not in this list, the method treats it as an error.
//   - responseObj: any - An optional parameter where the response body can be unmarshaled into. This is useful for directly obtaining a structured representation of the response data.
//
// Returns:
//   - *Response: The response from the HTTP request.
//   - error: An error if the request fails. Possible error types include:
//   - UrlFormatError: Returned if the URL is invalid or not absolute.
//   - UnexpectedHttpStatusCodeError: Returned if the response status code is not acceptable.
//
// If no scopes are provided, the method attempts to infer the scope from the URL. The URL is validated to ensure it is absolute and properly formatted.
// The HTTP request is then prepared and executed. The response status code is checked against the list of acceptable status codes. If the status code
// is not acceptable, an error is returned. If a responseObj is provided, the response body is unmarshaled into this object.
func (client *Client) Execute(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*Response, error) {
	if len(scopes) == 0 {
		// if no scopes are provided, try to guess the scope from the URL.
		scope, err := tryGetScopeFromURL(url, client.Config.Urls)
		if err != nil {
			return nil, err
		}

		scopes = []string{scope}
	}

	if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
		return nil, customerrors.NewUrlFormatError(url, e)
	}

	for {
		token, err := client.BaseAuth.GetTokenForScopes(ctx, scopes)

		// if method == "POST" {
		// 	// save token to file for debugging
		// 	err = os.WriteFile("token.txt", []byte(*token), 0644)
		// }

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
			return resp, fmt.Errorf("Error making %s request to %s. %w", request.Method, request.RequestURI, err)
		}

		isAcceptable := len(acceptableStatusCodes) > 0 && array.Contains(acceptableStatusCodes, resp.HttpResponse.StatusCode)
		if isAcceptable {
			if responseObj != nil && len(resp.BodyAsBytes) > 0 {
				err = resp.MarshallTo(responseObj)
				if err != nil {
					return resp, fmt.Errorf("Error marshalling response to json. %w", err)
				}
			}

			return resp, nil
		}

		isRetryable := array.Contains(retryableStatusCodes, resp.HttpResponse.StatusCode)
		if !isRetryable {
			return resp, customerrors.NewUnexpectedHttpStatusCodeError(acceptableStatusCodes, resp.HttpResponse.StatusCode, resp.HttpResponse.Status, resp.BodyAsBytes)
		}

		waitFor := retryAfter(ctx, resp.HttpResponse)

		tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", resp.HttpResponse.StatusCode, url, waitFor))

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return resp, err
		}
	}
}

// RetryAfterDefault returns a random duration between 10 and 20 seconds.
func DefaultRetryAfter() time.Duration {
	return time.Duration((rand.Intn(10) + 10)) * time.Second
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
	case strings.LastIndex(url, cloudConfig.PowerAppsAdvisor) != -1:
		return cloudConfig.PowerAppsAdvisorScope, nil
	case strings.LastIndex(url, cloudConfig.AdminPowerPlatformUrl) != -1:
		return constants.PPAC_SCOPE, nil
	case strings.LastIndex(url, "csanalytics") != -1:
		return cloudConfig.AnalyticsScope, nil
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
