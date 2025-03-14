// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package analytics_data_export

type AnalyticsDataResponse struct {
	Value []AnalyticsDataDto `json:"value"`
}

type AnalyticsDataDto struct {
	ID               string           `json:"id"`
	Source           string           `json:"source"`
	Environments     []EnvironmentDto `json:"environments"`
	Status           []StatusDto      `json:"status"`
	Sink             SinkDto          `json:"sink"`
	PackageName      string           `json:"packageName"`
	Scenarios        []string         `json:"scenarios"`
	ResourceProvider string           `json:"resourceProvider"`
	AiType           string           `json:"aiType"`
}

type StatusDto struct {
	Name      string  `json:"name"`
	State     string  `json:"state"`
	LastRunOn string  `json:"lastRunOn"`
	Message   *string `json:"message"`
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
	Type              string `json:"type"`
	SubscriptionId    string `json:"subscriptionId,omitempty"`
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	ResourceName      string `json:"resourceName"`
	Key               string `json:"key"`
}

// GatewayClusterDto represents a gateway cluster in the Power Platform.
type GatewayClusterDto struct {
	ClusterNumber   string `json:"clusterNumber"`
	GeoName         string `json:"geoName"`
	Environment     string `json:"environment"`
	ClusterType     string `json:"clusterType"`
	ClusterCategory string `json:"clusterCategory"`
	ClusterName     string `json:"clusterName"`
	GeoLongName     string `json:"geoLongName"`
}
