// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_checker_rules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

type Client struct {
	Api               *api.Client
	environmentClient environment.Client
}

func NewSolutionCheckerRulesClient(apiClient *api.Client) *Client {
	return &Client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

// Data transfer objects.
type RuleDto struct {
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

func (c *Client) GetSolutionCheckerRules(ctx context.Context, environmentId string) ([]RuleDto, error) {
	env, err := c.environmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment details for %s: %w", environmentId, err)
	}

	if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
		return nil, fmt.Errorf("could not find PowerAppsAdvisor endpoint for environment %s", environmentId)
	}

	// Parse and get host from PowerAppsAdvisor URL
	advisorURL, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PowerAppsAdvisor URL: %w", err)
	}

	// Now get the solution checker rules from the advisor endpoint with required query parameters
	rulesBaseUrl := fmt.Sprintf("https://%s/api/rule", advisorURL.Host)
	rulesUrl, err := url.Parse(rulesBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rules URL: %w", err)
	}

	// Add the required query parameters
	queryParams := url.Values{}
	queryParams.Add("ruleset", constants.SOLUTION_CHECKER_RULESET_ID)
	queryParams.Add("api-version", "2.0")
	rulesUrl.RawQuery = queryParams.Encode()

	var rules []RuleDto
	_, err = c.Api.Execute(ctx, nil, "GET", rulesUrl.String(), nil, nil, []int{http.StatusOK}, &rules)
	if err != nil {
		return nil, fmt.Errorf("failed to get solution checker rules: %w", err)
	}

	return rules, nil
}
