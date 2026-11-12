// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
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
	sleep     func(ctx context.Context, duration time.Duration) error
}

// getTenantIsolationPolicy retrieves the tenant isolation policy for the specified tenant.
// If no policy exists, it returns nil without error (HTTP 404).
//
// Parameters:
//   - ctx: Context for the request with cancellation and timeout capabilities
//   - tenantId: The ID of the tenant whose isolation policy should be retrieved
//
// Returns:
//   - *TenantIsolationPolicyDto: The retrieved policy or nil if not found
//   - error: Any error encountered during the API operation
func (client *Client) getTenantIsolationPolicy(ctx context.Context, tenantId string) (*TenantIsolationPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}

	var policy TenantIsolationPolicyDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &policy)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve tenant isolation policy for tenant %s: %w", tenantId, err)
	}

	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Info(ctx, fmt.Sprintf("Tenant isolation policy not found for tenant %s", tenantId))
		return nil, nil
	}

	return &policy, nil
}

// createOrUpdateTenantIsolationPolicy applies a new or updates an existing tenant isolation policy.
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
func (client *Client) createOrUpdateTenantIsolationPolicy(ctx context.Context, tenantId string, policy TenantIsolationPolicyDto) (*TenantIsolationPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}

	var updatedPolicy TenantIsolationPolicyDto
	resp, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, policy, []int{http.StatusOK, http.StatusAccepted}, &updatedPolicy)
	if err != nil {
		return nil, fmt.Errorf("could not create tenant isolation policy: %w", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusAccepted {
		// Handle async operation
		tflog.Info(ctx, fmt.Sprintf("Tenant isolation policy operation is asynchronous for tenant %s, waiting for completion...", tenantId))
		_, err := client.doWaitForLifecycleOperationStatus(ctx, resp)
		if err != nil {
			return nil, fmt.Errorf("error waiting for tenant isolation policy operation to complete: %w", err)
		}

		// Get fresh state after async operation
		finalPolicy, err := client.getTenantIsolationPolicy(ctx, tenantId)
		if err != nil {
			return nil, fmt.Errorf("error getting tenant isolation policy after async operation: %w", err)
		}

		return finalPolicy, nil
	}

	return &updatedPolicy, nil
}

// getRetryAfterDuration handles Retry-After header parsing.
// TODO: Consider using shared implementation from api package once exported.
func getRetryAfterDuration(resp *http.Response) time.Duration {
	// Default value if header not present or invalid
	defaultDuration := 5 * time.Second

	if resp == nil {
		return defaultDuration
	}

	// Check for Retry-After header
	retryAfter := resp.Header.Get(constants.HEADER_RETRY_AFTER)
	if retryAfter == "" {
		return defaultDuration
	}

	// Try to parse as seconds (integer)
	seconds, err := strconv.Atoi(retryAfter)
	if err == nil && seconds > 0 {
		duration := time.Duration(seconds) * time.Second

		// For safety, ensure we have a minimum wait time and cap the maximum
		if duration < 2*time.Second {
			duration = 2 * time.Second
		} else if duration > 60*time.Second {
			duration = 60 * time.Second
		}

		return duration
	}

	// If header value is not a valid integer, return default
	return defaultDuration
}

// doWaitForLifecycleOperationStatus polls an asynchronous operation until completion.
// It follows the location header from the initial response and waits the delay of the Retry-After
// header (default 5s, bounded to 2s-60s) before the first poll and between polls, including after a
// 409. The operation service can answer 404 OperationNotFound for an operation it has just accepted,
// so a 404 is polled again until the delays spent on 404 answers reach
// constants.TENANT_ISOLATION_OPERATION_NOT_FOUND_TIMEOUT.
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
	tflog.Debug(ctx, fmt.Sprintf("Location Header: %s", locationHeader))

	if locationHeader == "" {
		return nil, errors.New("no Location or Operation-Location header found in async response")
	}

	// Polling the moment the operation is accepted races its registration, so wait first.
	if err := client.sleepWithContext(ctx, getRetryAfterDuration(response.HttpResponse)); err != nil {
		return nil, fmt.Errorf("polling interrupted: %w", err)
	}

	notFoundWaited := time.Duration(0)
	for {
		apiResp, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted, http.StatusConflict, http.StatusNotFound}, nil)
		if err != nil {
			return nil, fmt.Errorf("error polling operation status: %w", err)
		}

		statusCode := apiResp.HttpResponse.StatusCode
		tflog.Debug(ctx, fmt.Sprintf("Operation poll status: %s", apiResp.HttpResponse.Status))

		// If we get a 200 OK, the operation is complete
		if statusCode == http.StatusOK {
			// Return nil since we'll get the fresh state later
			return nil, nil
		}

		// Get the next wait time from the Retry-After header if available
		waitTime := getRetryAfterDuration(apiResp.HttpResponse)

		if statusCode == http.StatusNotFound {
			if notFoundWaited >= constants.TENANT_ISOLATION_OPERATION_NOT_FOUND_TIMEOUT {
				statusErr := customerrors.NewUnexpectedHttpStatusCodeError([]int{http.StatusOK, http.StatusAccepted, http.StatusConflict}, statusCode, apiResp.HttpResponse.Status, apiResp.BodyAsBytes)
				return nil, fmt.Errorf("error polling operation status: the operation was still not found after waiting %s: %w", notFoundWaited, statusErr)
			}
			notFoundWaited += waitTime
			tflog.Debug(ctx, fmt.Sprintf("Operation %s is not found yet, %s of %s spent waiting for it", locationHeader, notFoundWaited, constants.TENANT_ISOLATION_OPERATION_NOT_FOUND_TIMEOUT))
		}

		tflog.Debug(ctx, fmt.Sprintf("Waiting for %s before polling again", waitTime))

		// Wait before polling again
		err = client.sleepWithContext(ctx, waitTime)
		if err != nil {
			return nil, fmt.Errorf("polling interrupted: %w", err)
		}
	}
}

// sleepWithContext waits between operation polls. Tests replace it through the sleep field to observe
// the delays, because the api client does not sleep in test mode.
func (client *Client) sleepWithContext(ctx context.Context, duration time.Duration) error {
	if client.sleep != nil {
		return client.sleep(ctx, duration)
	}
	return client.Api.SleepWithContext(ctx, duration)
}
