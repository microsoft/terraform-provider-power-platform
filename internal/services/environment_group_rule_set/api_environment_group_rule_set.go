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
)

func newEnvironmentGroupRuleSetClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string) (*environmentGroupRuleSetValueSetDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := environmentGroupRuleSetDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &environmentGroupRuleSet)
	if err != nil {
		return nil, err
	}

	if len(environmentGroupRuleSet.Value) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet.Value[0], nil
}

func (client *client) CreateEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string, newEnvironmentGroupRuleSet createEnvironmentGroupRuleSetDto) (*environmentGroupRuleSetValueSetDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := environmentGroupRuleSetDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentGroupRuleSet, []int{http.StatusCreated}, &newEnvironmentGroupRuleSet)
	if err != nil {
		return nil, err
	}

	if len(environmentGroupRuleSet.Value) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet.Value[0], nil
}
