// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_checker_rules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
)

type client interface {
	GetSolutionCheckerRules(ctx context.Context, environmentId string) ([]ruleDto, error)
}

type clientImpl struct {
	Api *api.Client
}

func newSolutionCheckerRulesClient(apiClient *api.Client) client {
	return &clientImpl{
		Api: apiClient,
	}
}

// Data transfer objects.
type ruleDto struct {
	Code            string `json:"code"`
	Description     string `json:"description"`
	Summary         string `json:"summary"`
	HowToFix        string `json:"howToFix"`
	GuidanceUrl     string `json:"guidanceUrl"`
	ComponentType   int    `json:"componentType"`
	PrimaryCategory int    `json:"primaryCategory"`
	Include         bool   `json:"include"`
	Severity        int    `json:"severity"`
}

func (c *clientImpl) GetSolutionCheckerRules(ctx context.Context, environmentId string) ([]ruleDto, error) {
	// Get the environment to find the PowerAppsAdvisor endpoint
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   c.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	var environmentResponse map[string]any
	_, err := c.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &environmentResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment details for %s: %w", environmentId, err)
	}

	// Extract the PowerAppsAdvisor endpoint
	var advisorEndpoint string
	if properties, ok := environmentResponse["properties"].(map[string]any); ok {
		if runtimeEndpoints, ok := properties["runtimeEndpoints"].(map[string]any); ok {
			if endpoint, ok := runtimeEndpoints["microsoft.PowerAppsAdvisor"].(string); ok {
				advisorEndpoint = endpoint
			}
		}
	}

	if advisorEndpoint == "" {
		return nil, fmt.Errorf("could not find PowerAppsAdvisor endpoint for environment %s", environmentId)
	}

	// Now get the solution checker rules from the advisor endpoint with required query parameters
	rulesBaseUrl := fmt.Sprintf("%s/api/rule", advisorEndpoint)
	rulesUrl, err := url.Parse(rulesBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rules URL: %w", err)
	}

	// Add the required query parameters
	queryParams := url.Values{}
	queryParams.Add("ruleset", "0ad12346-e108-40b8-a956-9a8f95ea18c9")
	queryParams.Add("api-version", "2.0")
	rulesUrl.RawQuery = queryParams.Encode()

	var rules []ruleDto
	_, err = c.Api.Execute(ctx, nil, "GET", rulesUrl.String(), nil, nil, []int{http.StatusOK}, &rules)
	if err != nil {
		return nil, fmt.Errorf("failed to get solution checker rules: %w", err)
	}

	return rules, nil
}
