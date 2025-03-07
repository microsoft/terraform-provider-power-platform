// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func NewTenantIsolationPolicyClient(apiClient *api.Client, tenantClient tenant.Client) Client {
	return Client{
		Api:       apiClient,
		TenantApi: tenantClient,
	}
}

type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

// GetTenantIsolationPolicy retrieves the tenant isolation policy for the specified tenant.
// If no policy exists, it returns nil without error (HTTP 404).
//
// Parameters:
//   - ctx: Context for the request with cancellation and timeout capabilities
//   - tenantId: The ID of the tenant whose isolation policy should be retrieved
//
// Returns:
//   - *TenantIsolationPolicyDto: The retrieved policy or nil if not found
//   - error: Any error encountered during the API operation
func (client *Client) GetTenantIsolationPolicy(ctx context.Context, tenantId string) (*TenantIsolationPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}

	var policy TenantIsolationPolicyDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &policy)
	if err != nil {
		return nil, err
	}

	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Info(ctx, fmt.Sprintf("Tenant isolation policy not found for tenant %s", tenantId))
		return nil, nil
	}

	return &policy, nil
}

// CreateOrUpdateTenantIsolationPolicy applies a new or updates an existing tenant isolation policy.
// This operation may be asynchronous (HTTP 202) or synchronous (HTTP 200).
// For asynchronous operations, it polls until completion before returning the result.
//
// Parameters:
//   - ctx: Context for the request with cancellation and timeout capabilities
//   - tenantId: The ID of the tenant whose policy should be created or updated
//   - policy: The isolation policy configuration to apply
//
// Returns:
//   - *TenantIsolationPolicyDto: The created or updated policy
//   - error: Any error from the operation or polling process
func (client *Client) CreateOrUpdateTenantIsolationPolicy(ctx context.Context, tenantId string, policy TenantIsolationPolicyDto) (*TenantIsolationPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}

	var updatedPolicy TenantIsolationPolicyDto
	resp, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, policy, []int{http.StatusOK, http.StatusAccepted}, &updatedPolicy)
	if err != nil {
		return nil, fmt.Errorf("Could not create tenant isolation policy: %v", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusAccepted {
		// Handle async operation
		tflog.Info(ctx, fmt.Sprintf("Tenant isolation policy operation is asynchronous for tenant %s, waiting for completion...", tenantId))
		_, err := client.doWaitForLifecycleOperationStatus(ctx, resp)
		if err != nil {
			return nil, fmt.Errorf("Error waiting for tenant isolation policy operation to complete: %v", err)
		}

		// Get fresh state after async operation
		finalPolicy, err := client.GetTenantIsolationPolicy(ctx, tenantId)
		if err != nil {
			return nil, fmt.Errorf("Error getting tenant isolation policy after async operation: %v", err)
		}

		return finalPolicy, nil
	}

	return &updatedPolicy, nil
}

// getRetryAfterDuration parses the Retry-After header from an HTTP response.
// It converts the header value to a time.Duration, defaulting to 5 seconds
// if the header is missing or invalid.
func getRetryAfterDuration(resp *http.Response) time.Duration {
	// Default value if header not present or invalid
	defaultDuration := 5 * time.Second

	if resp == nil {
		return defaultDuration
	}

	// Check for Retry-After header
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return defaultDuration
	}

	// Try to parse as seconds (integer)
	seconds, err := strconv.Atoi(retryAfter)
	if err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	// If header value is not a valid integer, return default
	return defaultDuration
}

// doWaitForLifecycleOperationStatus polls an asynchronous operation until completion.
// It follows the location header from the initial response, respects Retry-After
// headers to control polling frequency, and implements exponential backoff with
// minimum 2s and maximum 60s between attempts.
//
// Parameters:
//   - ctx: Context for the request with cancellation and timeout capabilities
//   - response: The initial API response containing operation location headers
//
// Returns:
//   - *TenantIsolationPolicyDto: Always nil (caller should get fresh state)
//   - error: Any error encountered during the polling process
func (client *Client) doWaitForLifecycleOperationStatus(ctx context.Context, response *api.Response) (*TenantIsolationPolicyDto, error) {
	locationHeader := response.GetHeader(constants.HEADER_LOCATION)
	if locationHeader == "" {
		locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
	}
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	if locationHeader == "" {
		return nil, fmt.Errorf("no Location or Operation-Location header found in async response")
	}

	for {
		apiResp, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, nil)
		if err != nil {
			return nil, fmt.Errorf("error polling operation status: %v", err)
		}

		tflog.Debug(ctx, fmt.Sprintf("Operation poll status: %s", apiResp.HttpResponse.Status))

		// If we get a 200 OK, the operation is complete
		if apiResp.HttpResponse.StatusCode == http.StatusOK {
			// Return nil since we'll get the fresh state later
			return nil, nil
		}

		// Get the next wait time from the Retry-After header if available
		waitTime := getRetryAfterDuration(apiResp.HttpResponse)

		// For safety, ensure we have a minimum wait time and cap the maximum
		if waitTime < 2*time.Second {
			waitTime = 2 * time.Second
		} else if waitTime > 60*time.Second {
			waitTime = 60 * time.Second
		}

		tflog.Debug(ctx, fmt.Sprintf("Waiting for %s before polling again", waitTime))

		// Wait before polling again
		err = client.Api.SleepWithContext(ctx, waitTime)
		if err != nil {
			return nil, fmt.Errorf("polling interrupted: %v", err)
		}
	}
}
