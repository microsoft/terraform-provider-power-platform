// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

type EnvironmentTemplateItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Properties struct {
		DisplayName    string `json:"displayName"`
		IsDisabled     bool   `json:"isDisabled"`
		DisabledReason struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"disabledReason"`
		IsCustomerEngagement         bool `json:"isCustomerEngagement"`
		IsSupportedForResetOperation bool `json:"isSupportedForResetOperation"`
	} `json:"properties"`
}

type EnvironmentTemplatesDto struct {
	Standard               []EnvironmentTemplateItem `json:"standard"`
	Premium                []EnvironmentTemplateItem `json:"premium"`
	Developer              []EnvironmentTemplateItem `json:"developer"`
	Basic                  []EnvironmentTemplateItem `json:"basic"`
	Production             []EnvironmentTemplateItem `json:"production"`
	Sandbox                []EnvironmentTemplateItem `json:"sandbox"`
	Trial                  []EnvironmentTemplateItem `json:"trial"`
	Default                []EnvironmentTemplateItem `json:"default"`
	Support                []EnvironmentTemplateItem `json:"support"`
	SubscriptionBasedTrial []EnvironmentTemplateItem `json:"subscriptionBasedTrial"`
	Teams                  []EnvironmentTemplateItem `json:"teams"`
	Platform               []EnvironmentTemplateItem `json:"platform"`
}
