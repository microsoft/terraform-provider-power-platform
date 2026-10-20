// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	499,                            // 499 is retryable because the client may have closed the connection.
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
			return strings.Contains(wwwAuthenticate, constants.CAE_CHALLENGE_CLAIMS_INDICATOR) &&
				strings.Contains(wwwAuthenticate, constants.CAE_CHALLENGE_INSUFFICIENT_CLAIMS_INDICATOR)
		}
	}

	return false
}

// Execute executes an HTTP request with the given method, url, headers, and body.
//
// Parameters:
//   - ctx: context.Context - Provides context for the request, allowing for timeout and cancellation control.
//   - scopes: []string - A list of scopes that the request should be associated with. If no scopes are provided, the method infers the scope from the host of the URL, which must be an https URL of a known Power Platform or Dataverse host.
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
//   - An error is also returned, before any request is sent, if no scopes are provided and none can be inferred from the URL.
//   - UnexpectedHttpStatusCodeError: Returned if the response status code is not acceptable.
//
// If no scopes are provided, the method attempts to infer the scope from the URL. The URL is validated to ensure it is absolute and properly formatted.
// The HTTP request is then prepared and executed. The response status code is checked against the list of acceptable status codes. If the status code
// is not acceptable, an error is returned. If a responseObj is provided, the response body is unmarshaled into this object.
// Transient statuses are retried at most constants.MAX_RETRY_COUNT times (401 at most constants.MAX_UNAUTHORIZED_RETRIES times),
// after which the last response is returned together with an error wrapping UnexpectedHttpStatusCodeError.
func (client *Client) Execute(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*Response, error) {
	return client.doExecute(ctx, scopes, method, url, headers, body, acceptableStatusCodes, responseObj, retryTransientStatuses)
}

// ExecuteWithoutRetry executes exactly one HTTP request. It is intended for
// non-idempotent mutation starts whose outcome can be ambiguous when the
// service returns a transient status after accepting the request. Callers must
// recover the remote operation by its durable identity rather than replaying
// the request.
func (client *Client) ExecuteWithoutRetry(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*Response, error) {
	return client.doExecute(ctx, scopes, method, url, headers, body, acceptableStatusCodes, responseObj, retryNever)
}

// retryBehavior selects whether doExecute replays requests that fail with a
// transient (retryable) HTTP status. Non-idempotent mutation starts must never
// be replayed, so they run with retryNever.
type retryBehavior int

const (
	retryTransientStatuses retryBehavior = iota
	retryNever
)

func (client *Client) doExecute(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any, retry retryBehavior) (*Response, error) {
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

	var lastRetryableResponse *Response
	retries := 0
	unauthorizedRetries := 0
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
			// If the operation deadline/cancellation stops a retry loop, retain the last real HTTP
			// failure. doRequest synthesizes an empty 503 for a transport cancellation, which would
			// otherwise erase the service status/body that caused the retries.
			if lastRetryableResponse != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				return lastRetryableResponse, err
			}
			return resp, err
		}

		err = validateNoManagementApplicationPermissionsForBapiRequest(resp)
		if err != nil {
			return resp, err
		}

		isAcceptable := len(acceptableStatusCodes) > 0 && helpers.ArrayContains(acceptableStatusCodes, resp.HttpResponse.StatusCode)
		if isAcceptable {
			if responseObj != nil && len(resp.BodyAsBytes) > 0 {
				err = resp.MarshallTo(responseObj)
				if err != nil {
					// The status was accepted, so the operation happened; only its result was lost.
					return resp, RequestSentError{Err: fmt.Errorf("Error marshalling response to json. %w", err)}
				}
			}

			return resp, nil
		}

		statusCode := resp.HttpResponse.StatusCode
		if err := checkRetryAllowed(retry, acceptableStatusCodes, resp, retries, unauthorizedRetries); err != nil {
			return resp, err
		}
		if statusCode == http.StatusUnauthorized {
			unauthorizedRetries++
		}
		retries++
		lastRetryableResponse = resp

		waitFor := retryAfter(ctx, resp.HttpResponse)

		tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", resp.HttpResponse.StatusCode, url, waitFor))

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return resp, err
		}
	}
}

// checkRetryAllowed returns nil when a response with an unacceptable status may be retried, or the
// error to surface when the status is permanent or its retry budget is exhausted.
func checkRetryAllowed(retry retryBehavior, acceptableStatusCodes []int, resp *Response, retries, unauthorizedRetries int) error {
	statusCode := resp.HttpResponse.StatusCode
	statusErr := customerrors.NewUnexpectedHttpStatusCodeError(acceptableStatusCodes, statusCode, resp.HttpResponse.Status, resp.BodyAsBytes)

	if statusCode == http.StatusUnauthorized {
		// The credential was refused, so the request never reached the operation and replaying
		// it with a freshly acquired token cannot duplicate a non-idempotent mutation. A 401 caused
		// by missing permissions is permanent, so the replays are bounded in every retry mode.
		if unauthorizedRetries >= constants.MAX_UNAUTHORIZED_RETRIES {
			return fmt.Errorf("request was still unauthorized after %d retries with a fresh token: %w", unauthorizedRetries, statusErr)
		}
	} else if retry == retryNever || !helpers.ArrayContains(retryableStatusCodes, statusCode) {
		return statusErr
	}

	if retries >= constants.MAX_RETRY_COUNT {
		return fmt.Errorf("maximum retries (%d) reached: %w", constants.MAX_RETRY_COUNT, statusErr)
	}
	return nil
}

func (client *Client) HandleNotFoundResponse(resp *Response) error {
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found at '%s'", resp.HttpResponse.Request.URL)
	}
	return nil
}

func (client *Client) HandleForbiddenResponse(resp *Response) error {
	if resp.HttpResponse.StatusCode == http.StatusForbidden {
		return customerrors.NewAccessDeniedError(requestUrl(resp), resp.BodyAsBytes)
	}
	return nil
}

func requestUrl(resp *Response) string {
	if resp.HttpResponse.Request == nil || resp.HttpResponse.Request.URL == nil {
		return ""
	}
	return resp.HttpResponse.Request.URL.String()
}

func validateNoManagementApplicationPermissionsForBapiRequest(resp *Response) error {
	if resp.HttpResponse.StatusCode == http.StatusForbidden && len(resp.BodyAsBytes) > 0 && strings.Contains(string(resp.BodyAsBytes), "does not have permission to access the path") {
		return errors.New(constants.NO_MANAGEMENT_APPLICATION_ERROR_MSG)
	}
	return nil
}

// RetryAfterDefault returns a random duration between 10 and 20 seconds.
func DefaultRetryAfter() time.Duration {
	return time.Duration((rand.Intn(10) + 10)) * time.Second
}

// SleepWithContext sleeps for the given duration or until the context is canceled.
func (client *Client) SleepWithContext(ctx context.Context, duration time.Duration) error {
	// Cancellation is part of the API contract even in test mode. Check it before bypassing the
	// physical sleep so retry-loop tests cannot spin forever after their deadline has elapsed.
	if err := ctx.Err(); err != nil {
		return err
	}
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

// dataverseDomains are the domains of Dataverse environment hosts in every supported cloud.
var dataverseDomains = []string{
	constants.PUBLIC_DATAVERSE_DOMAIN,
	constants.USGOVHIGH_DATAVERSE_DOMAIN,
	constants.USDOD_DATAVERSE_DOMAIN,
	constants.CHINA_DATAVERSE_DOMAIN,
	constants.EX_DATAVERSE_DOMAIN,
	constants.RX_DATAVERSE_DOMAIN,
}

// analyticsDomains are the domains of the admin analytics hosts, which are chosen by geography rather than by cloud.
var analyticsDomains = []string{
	constants.PUBLIC_ANALYTICS_DOMAIN,
	constants.USGOV_ANALYTICS_DOMAIN,
	constants.USDOD_ANALYTICS_DOMAIN,
}

// tryGetScopeFromURL returns the authorization scope for a request to the given API URL.
//
// The scope is chosen from the parsed host of the URL alone. A host belongs to a domain only when it is the
// domain itself or a subdomain of it, so a URL that merely mentions a known domain elsewhere (in its path,
// query, fragment, user info, or as the prefix of a longer host name) never gets that domain's token.
// Only https URLs get a scope, and a host that is neither a known API host nor a Dataverse host is rejected
// instead of being sent a token minted for itself.
func tryGetScopeFromURL(rawUrl string, cloudConfig config.ProviderConfigUrls) (string, error) {
	u, err := neturl.Parse(rawUrl)
	if err != nil || !u.IsAbs() {
		return "", customerrors.NewUrlFormatError(rawUrl, err)
	}
	if !strings.EqualFold(u.Scheme, constants.HTTPS) {
		return "", fmt.Errorf("refusing to authorize a request to '%s': only %s urls can be authorized", rawUrl, constants.HTTPS)
	}

	host := strings.ToLower(u.Hostname())
	switch {
	case isHostInDomains(host, cloudConfig.BapiUrl, cloudConfig.PowerAppsUrl):
		return cloudConfig.PowerAppsScope, nil
	case isHostInDomains(host, cloudConfig.PowerPlatformUrl):
		return cloudConfig.PowerPlatformScope, nil
	case isHostInDomains(host, cloudConfig.PowerAppsAdvisor):
		return cloudConfig.PowerAppsAdvisorScope, nil
	case isHostInDomains(host, cloudConfig.AdminPowerPlatformUrl):
		return constants.PPAC_SCOPE, nil
	case isHostInDomains(host, analyticsDomains...):
		return cloudConfig.AnalyticsScope, nil
	case isHostInDomains(host, cloudConfig.LicensingUrl),
		isHostInDomains(host, dataverseDomains...):
		return constants.HTTPS + "://" + u.Host + "/.default", nil
	default:
		return "", fmt.Errorf("refusing to authorize a request to '%s': '%s' is not a known Power Platform or Dataverse host, so no authorization scope can be inferred for it", rawUrl, host)
	}
}

// isHostInDomains reports whether the host is one of the domains or a subdomain of one of them.
func isHostInDomains(host string, domains ...string) bool {
	if host == "" {
		return false
	}
	for _, domain := range domains {
		domain = strings.ToLower(domain)
		if domain != "" && (host == domain || strings.HasSuffix(host, "."+domain)) {
			return true
		}
	}
	return false
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
