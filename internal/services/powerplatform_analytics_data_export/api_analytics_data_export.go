// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

import (
	"context"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func NewTelemetryExport(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) GetTelemetryExport(ctx context.Context) (*TelemetryExportDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "/api/v2/connections",
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	telemetrylinks := TelemetryExportDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &telemetrylinks)

	return &telemetrylinks, err
}

func (client *Client) CreateTelemetryExport(ctx context.Context, telemetryToCreate TelemetryExportCreateDto) (*TelemetryExportDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "/api/v2/sinks/appinsights/connections",
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	telemetrylinks := &TelemetryExportDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, telemetryToCreate, []int{http.StatusCreated}, telemetrylinks)
	if err != nil {
		return nil, err
	}

	return telemetrylinks, err
}
