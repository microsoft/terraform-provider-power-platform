// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

type itemDto struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Location   string            `json:"location"`
	Properties propertiesItemDto `json:"properties"`
}

type propertiesItemDto struct {
	DisplayName                  string                `json:"displayName"`
	IsDisabled                   bool                  `json:"isDisabled"`
	DisabledReason               itemDisabledReasonDto `json:"disabledReason"`
	IsCustomerEngagement         bool                  `json:"isCustomerEngagement"`
	IsSupportedForResetOperation bool                  `json:"isSupportedForResetOperation"`
}

type itemDisabledReasonDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type environmentTemplateDto struct {
	Standard               []itemDto `json:"standard"`
	Premium                []itemDto `json:"premium"`
	Developer              []itemDto `json:"developer"`
	Basic                  []itemDto `json:"basic"`
	Production             []itemDto `json:"production"`
	Sandbox                []itemDto `json:"sandbox"`
	Trial                  []itemDto `json:"trial"`
	Default                []itemDto `json:"default"`
	Support                []itemDto `json:"support"`
	SubscriptionBasedTrial []itemDto `json:"subscriptionBasedTrial"`
	Teams                  []itemDto `json:"teams"`
	Platform               []itemDto `json:"platform"`
}
