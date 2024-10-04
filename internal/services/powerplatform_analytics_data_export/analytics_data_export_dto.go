// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

type TelemetryExportDto struct {
	Id               string        `json:"id"`
	Source           string        `json:"source"`
	Environments     []Environment `json:"environments"`
	Sink             Sink          `json:"sink"`
	Statuses         []Status      `json:"status"`
	PackageName      string        `json:"packageName"`
	Scenarios        []string      `json:"scenarios"`
	ResourceProvider string        `json:"resourceProvider"`
	AiType           string        `json:"aiType"`
}

type TelemetryExportCreateDto struct {
	Source           string        `json:"source"`
	Environments     []Environment `json:"environments"`
	Sink             Sink          `json:"sink"`
	PackageName      string        `json:"packageName"`
	Scenarios        []string      `json:"scenarios"`
	ResourceProvider string        `json:"resourceProvider"`
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
