// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

type AnalyticsDataDto struct {
	Source           string           `json:"source"`
	Environments     []EnvironmentDto `json:"environments"`
	Sink             SinkDto          `json:"sink"`
	Status           []Status         `json:"status"`
	PackageName      string           `json:"packageName"`
	Scenarios        []string         `json:"scenarios"`
	ResourceProvider string           `json:"resourceProvider"`
	AiType           string           `json:"aiType"`
}

type AnalyticsDataCreateDto struct {
	Source           string           `json:"source"`
	Environments     []EnvironmentDto `json:"environments"`
	Sink             SinkDto          `json:"sink"`
	PackageName      string           `json:"packageName"`
	Scenarios        []string         `json:"scenarios"`
	ResourceProvider string           `json:"resourceProvider"`
}

type EnvironmentDto struct {
	EnvironmentId  string `json:"environmentId"`
	OrganizationId string `json:"organizationId"`
}

type SinkDto struct {
	ID                string `json:"id"`
	SubscriptionId    string `json:"subscriptionId"`
	ResourceGroupName string `json:"resourceGroupName"`
	ResourceName      string `json:"resourceName"`
	Key               string `json:"key"`
}
