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

// NewTenantIsolationPolicyClient creates a new client for tenant isolation policy operations.
func NewTenantIsolationPolicyClient(apiClient *api.Client, tenantClient tenant.Client) Client {
	return Client{
		Api:       apiClient,
		TenantApi: tenantClient,
	}
}

// Client represents the tenant isolation policy client.
type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

// GetTenantIsolationPolicy retrieves the tenant isolation policy.
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

// CreateOrUpdateTenantIsolationPolicy creates or updates the tenant isolation policy.
func (client *Client) CreateOrUpdateTenantIsolationPolicy(ctx context.Context, tenantId string, policy TenantIsolationPolicyDto) (*TenantIsolationPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", tenantId),
	}

	var updatedPolicy TenantIsolationPolicyDto
	resp, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, policy, []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, &updatedPolicy)
	if err != nil {
		return nil, fmt.Errorf("Could not create tenant isolation policy: %v", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusAccepted {
		// Handle async operation
		updatedPolicy, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
		if err != nil {
			return nil, fmt.Errorf("Error waiting for tenant isolation policy operation to complete: %v", err)
		}

		return updatedPolicy, nil
	}

	return &updatedPolicy, nil
}

// getRetryAfterDuration extracts the Retry-After header value and converts it to a duration
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

// DoWaitForLifecycleOperationStatus polls the operation status endpoint until the operation completes
// and returns the final tenant isolation policy resource from the response.
func (client *Client) DoWaitForLifecycleOperationStatus(ctx context.Context, response *api.Response) (*TenantIsolationPolicyDto, error) {
	locationHeader := response.GetHeader(constants.HEADER_LOCATION)
	if locationHeader == "" {
		locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
	}
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	if locationHeader == "" {
		return nil, fmt.Errorf("no Location or Operation-Location header found in async response")
	}

	// Initial wait time
	waitTime := getRetryAfterDuration(response.HttpResponse)

	for {
		// Use the TenantIsolationPolicyDto type directly to unmarshal the response
		var policyResponse TenantIsolationPolicyDto
		apiResp, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, &policyResponse)
		if err != nil {
			return nil, fmt.Errorf("error polling operation status: %v", err)
		}

		tflog.Debug(ctx, fmt.Sprintf("Operation poll status: %s", apiResp.HttpResponse.Status))

		// If we get a 200 OK, the operation is complete and the response contains the resource
		if apiResp.HttpResponse.StatusCode == http.StatusOK {
			// Return the unmarshalled policy directly from the response
			return &policyResponse, nil
		}

		// Get the next wait time from the Retry-After header if available
		waitTime = getRetryAfterDuration(apiResp.HttpResponse)

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
