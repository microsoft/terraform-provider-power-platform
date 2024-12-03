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

func (client *Client) GetAnalyticsDataExport(ctx context.Context) (*AnalyticsDataDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "api/v2/connections",
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	analyticdatalinks := AnalyticsDataDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &analyticdatalinks)

	return &analyticdatalinks, err
}

func (client *Client) CreateAnalyticsDataExport(ctx context.Context, analyticsdataToCreate AnalyticsDataCreateDto) (*AnalyticsDataDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "api/v2/sinks/appinsights/connections",
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	analyticdatalinks := &AnalyticsDataDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, analyticsdataToCreate, []int{http.StatusCreated}, analyticdatalinks)
	if err != nil {
		return nil, err
	}

	return analyticdatalinks, err
}

func (client *Client) UpdateAnalyticsDataExport(ctx context.Context, id string, analyticsdataToUpdate AnalyticsDataCreateDto) (*AnalyticsDataDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "api/v2/sinks/appinsights/connections/" + id,
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	analyticdatalinks := &AnalyticsDataDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, analyticsdataToUpdate, []int{http.StatusOK}, analyticdatalinks)
	if err != nil {
		return nil, err
	}

	return analyticdatalinks, err
}

func (client *Client) DeleteAnalyticsDataExport(ctx context.Context, id string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "api/v2/sinks/appinsights/connections/" + id,
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
	if err != nil {
		return err
	}

	return nil
}
