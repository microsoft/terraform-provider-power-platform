// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

import (
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type AnalyticsDataModel struct {
	Source           string        `tfsdk:"source"`
	Environments     []Environment `tfsdk:"environments"`
	Sink             Sink          `tfsdk:"sink"`
	PackageName      string        `tfsdk:"packageName"`
	Scenarios        []string      `tfsdk:"scenarios"`
	ResourceProvider string        `tfsdk:"resourceProvider"`
}

type Environment struct {
	EnvironmentId  string `tfsdk:"environmentId"`
	OrganizationId string `tfsdk:"organizationId"`
}

type Sink struct {
	ID                string `tfsdk:"id"`
	SubscriptionId    string `tfsdk:"subscriptionId"`
	ResourceGroupName string `tfsdk:"resourceGroupName"`
	ResourceName      string `tfsdk:"resourceName"`
	Key               string `tfsdk:"key"`
}

type Status struct {
	Name      string    `tfsdk:"name"`
	State     string    `tfsdk:"state"`
	LastRunOn time.Time `tfsdk:"lastRunOn"`
	Message   string    `tfsdk:"message"`
}

type AnalyticsExportDatasourceModel struct {
	Source           string        `tfsdk:"source"`
	Environments     []Environment `tfsdk:"environments"`
	Sink             Sink          `tfsdk:"sink"`
	PackageName      string        `tfsdk:"packageName"`
	Scenarios        []string      `tfsdk:"scenarios"`
	ResourceProvider string        `tfsdk:"resourceProvider"`
	Status           []Status      `tfsdk:"status"`
	AiType           string        `tfsdk:"aiType"`
}

type AnalyticsExportDataSource struct {
	helpers.TypeInfo
	AnalyticsExportData Client
}

type ResourceAnalyticsDataExport struct {
	helpers.TypeInfo
	ResourceAnalyticsDataExport Client
}
