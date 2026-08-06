// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_based_policy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

const apiVersion = "2024-10-01"

func NewRuleBasedPolicyClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) CreatePolicy(ctx context.Context, request ruleBasedPolicyRequestDto) (*ruleBasedPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/governance/ruleBasedPolicies",
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	policy := ruleBasedPolicyDto{}
	resp, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, request, []int{http.StatusOK, http.StatusCreated, http.StatusConflict}, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule-based policy: %w", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusConflict {
		return nil, &PolicyConflictError{}
	}

	return &policy, nil
}

// PolicyConflictError indicates that a policy already exists for the tenant.
type PolicyConflictError struct{}

func (e *PolicyConflictError) Error() string {
	return "a rule-based policy already exists for this tenant"
}

func (client *Client) GetPolicy(ctx context.Context, policyId string) (*ruleBasedPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/governance/ruleBasedPolicies/%s", policyId),
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	policy := ruleBasedPolicyDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to get rule-based policy: %w", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(
			fmt.Errorf("rule-based policy '%s' not found", policyId),
			constants.ERROR_OBJECT_NOT_FOUND,
			fmt.Sprintf("rule-based policy '%s' not found", policyId),
		)
	}

	return &policy, nil
}

func (client *Client) UpdatePolicy(ctx context.Context, policyId string, request ruleBasedPolicyRequestDto) (*ruleBasedPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/governance/ruleBasedPolicies/%s", policyId),
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	policy := ruleBasedPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, request, []int{http.StatusOK}, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule-based policy: %w", err)
	}

	return &policy, nil
}

func (client *Client) RemoveRuleFromPolicy(ctx context.Context, policyId string, policyName string, ruleSetId string, version string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/governance/ruleBasedPolicies/%s/removeRule", policyId),
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	request := removeRuleRequestDto{
		Name: policyName,
		RuleSets: []removeRuleSetDto{
			{
				Id:      ruleSetId,
				Version: version,
			},
		},
	}

	_, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, request, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return fmt.Errorf("failed to remove rule %s from policy %s: %w", ruleSetId, policyId, err)
	}

	return nil
}

func (client *Client) CreateEnvironmentGroupAssignment(ctx context.Context, policyId string, groupId string) (*ruleAssignmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/governance/ruleBasedPolicies/%s/environmentGroups/%s/assignments", policyId, groupId),
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	request := policyAssignmentRequestDto{
		AssignmentOverrides: []interface{}{},
	}

	assignment := ruleAssignmentDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, request, []int{http.StatusOK, http.StatusCreated}, &assignment)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment group assignment: %w", err)
	}

	return &assignment, nil
}

func (client *Client) ListAssignmentsByEnvironmentGroup(ctx context.Context, groupId string) ([]ruleAssignmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/governance/ruleBasedPolicies/environmentGroups/%s/assignments", groupId),
	}

	values := url.Values{}
	values.Add("api-version", apiVersion)
	values.Add("includeRuleSetCounts", "true")
	apiUrl.RawQuery = values.Encode()

	response := ruleAssignmentsResponseDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignments by environment group: %w", err)
	}

	return response.Value, nil
}
