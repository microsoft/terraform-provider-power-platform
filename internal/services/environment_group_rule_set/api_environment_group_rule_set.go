// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func NewEnvironmentGroupRuleSetClient(apiClient *api.Client, tenantClient tenant.Client) Client {
	return Client{
		Api:       apiClient,
		TenantApi: tenantClient,
	}
}

type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

func (client *Client) GetEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := environmentGroupRuleSetDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, &environmentGroupRuleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment group rule set: %w", err)
	}

	if resp.HttpResponse.StatusCode == http.StatusNoContent {
		return nil, customerrors.WrapIntoProviderError(
			fmt.Errorf("rule set '%s' not found", environmentGroupId),
			constants.ERROR_OBJECT_NOT_FOUND,
			fmt.Sprintf("rule set '%s' not found", environmentGroupId),
		)
	}

	if len(environmentGroupRuleSet.Value) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet.Value[0], nil
}

func (client *Client) CreateEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string, newEnvironmentGroupRuleSet EnvironmentGroupRuleSetValueSetDto) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := EnvironmentGroupRuleSetValueSetDto{}
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, newEnvironmentGroupRuleSet, []int{http.StatusCreated}, &environmentGroupRuleSet)

	if err != nil {
		return nil, fmt.Errorf("failed to create environment group rule set: %w", err)
	}

	// If empty Parameters is a valid API response, remove the following check.
	// Otherwise, consider clarifying its necessity and error messaging.
	if len(environmentGroupRuleSet.Parameters) == 0 {
		return nil, fmt.Errorf("no environment group ruleset parameters found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet, nil
}

func (client *Client) UpdateEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string, newEnvironmentGroupRuleSet EnvironmentGroupRuleSetValueSetDto) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/ruleSets/%s", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := EnvironmentGroupRuleSetValueSetDto{}
	_, err = client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, newEnvironmentGroupRuleSet, []int{http.StatusOK}, &environmentGroupRuleSet)

	if err != nil {
		return nil, fmt.Errorf("failed to update environment group rule set: %w", err)
	}

	return &environmentGroupRuleSet, nil
}

func (client *Client) DeleteEnvironmentGroupRuleSet(ctx context.Context, ruleSetId string) error {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/ruleSets/%s", ruleSetId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete environment group rule set: %w", err)
	}

	return nil
}

// policyUrl builds a rule-based policies URL. These endpoints live on the flat Power Platform
// host rather than the tenant-specific host used by the legacy rule sets endpoints.
func (client *Client) policyUrl(path string, extraQuery url.Values) string {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}

	values := url.Values{}
	values.Add("api-version", POLICY_API_VERSION)
	for k, vs := range extraQuery {
		for _, v := range vs {
			values.Add(k, v)
		}
	}
	apiUrl.RawQuery = values.Encode()

	return apiUrl.String()
}

func (client *Client) createPolicy(ctx context.Context, request ruleBasedPolicyRequestDto) (*ruleBasedPolicyDto, error) {
	policy := ruleBasedPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", client.policyUrl("/governance/ruleBasedPolicies", nil), nil, request, []int{http.StatusOK, http.StatusCreated}, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule-based policy: %w", err)
	}

	return &policy, nil
}

func (client *Client) getPolicy(ctx context.Context, policyId string) (*ruleBasedPolicyDto, error) {
	policy := ruleBasedPolicyDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", client.policyUrl(fmt.Sprintf("/governance/ruleBasedPolicies/%s", policyId), nil), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &policy)
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

func (client *Client) updatePolicy(ctx context.Context, policyId string, request ruleBasedPolicyRequestDto) (*ruleBasedPolicyDto, error) {
	policy := ruleBasedPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", client.policyUrl(fmt.Sprintf("/governance/ruleBasedPolicies/%s", policyId), nil), nil, request, []int{http.StatusOK}, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule-based policy: %w", err)
	}

	return &policy, nil
}

func (client *Client) deletePolicy(ctx context.Context, policyId string) error {
	_, err := client.Api.Execute(ctx, nil, "DELETE", client.policyUrl(fmt.Sprintf("/governance/ruleBasedPolicies/%s", policyId), nil), nil, nil, []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete rule-based policy: %w", err)
	}

	return nil
}

func (client *Client) deleteEnvironmentGroupAssignment(ctx context.Context, policyId, groupId string) error {
	path := fmt.Sprintf("/governance/ruleBasedPolicies/%s/environmentGroups/%s/assignments", policyId, groupId)
	_, err := client.Api.Execute(ctx, nil, "DELETE", client.policyUrl(path, nil), nil, nil, []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete environment group assignment: %w", err)
	}

	return nil
}

func (client *Client) removeRuleFromPolicy(ctx context.Context, policyId, policyName, ruleSetId, version string) error {
	request := removeRuleRequestDto{
		Name:     policyName,
		RuleSets: []removeRuleSetDto{{Id: ruleSetId, Version: version}},
	}

	_, err := client.Api.Execute(ctx, nil, "PATCH", client.policyUrl(fmt.Sprintf("/governance/ruleBasedPolicies/%s/removeRule", policyId), nil), nil, request, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return fmt.Errorf("failed to remove rule %s from policy %s: %w", ruleSetId, policyId, err)
	}

	return nil
}

func (client *Client) createEnvironmentGroupAssignment(ctx context.Context, policyId, groupId string) (*ruleAssignmentDto, error) {
	assignment := ruleAssignmentDto{}
	path := fmt.Sprintf("/governance/ruleBasedPolicies/%s/environmentGroups/%s/assignments", policyId, groupId)
	_, err := client.Api.Execute(ctx, nil, "POST", client.policyUrl(path, nil), nil, policyAssignmentRequestDto{AssignmentOverrides: []any{}}, []int{http.StatusOK, http.StatusCreated}, &assignment)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment group assignment: %w", err)
	}

	return &assignment, nil
}

func (client *Client) listAssignmentsByEnvironmentGroup(ctx context.Context, groupId string) ([]ruleAssignmentDto, error) {
	response := ruleAssignmentsResponseDto{}
	path := fmt.Sprintf("/governance/ruleBasedPolicies/environmentGroups/%s/assignments", groupId)
	_, err := client.Api.Execute(ctx, nil, "GET", client.policyUrl(path, url.Values{"includeRuleSetCounts": []string{"true"}}), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignments by environment group: %w", err)
	}

	return response.Value, nil
}
