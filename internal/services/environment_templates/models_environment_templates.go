// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

type ItemDto struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	Properties PropertiesItemDto `json:"properties"`
}

type PropertiesItemDto struct {
	DisplayName                  string                `json:"displayName"`
	IsDisabled                   bool                  `json:"isDisabled"`
	DisabledReason               ItemDisabledReasonDto `json:"disabledReason"`
	IsCustomerEngagement         bool                  `json:"isCustomerEngagement"`
	IsSupportedForResetOperation bool                  `json:"isSupportedForResetOperation"`
}

type ItemDisabledReasonDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Dto struct {
	Standard               []ItemDto `json:"standard"`
	Premium                []ItemDto `json:"premium"`
	Developer              []ItemDto `json:"developer"`
	Basic                  []ItemDto `json:"basic"`
	Production             []ItemDto `json:"production"`
	Sandbox                []ItemDto `json:"sandbox"`
	Trial                  []ItemDto `json:"trial"`
	Default                []ItemDto `json:"default"`
	Support                []ItemDto `json:"support"`
	SubscriptionBasedTrial []ItemDto `json:"subscriptionBasedTrial"`
	Teams                  []ItemDto `json:"teams"`
	Platform               []ItemDto `json:"platform"`
}
