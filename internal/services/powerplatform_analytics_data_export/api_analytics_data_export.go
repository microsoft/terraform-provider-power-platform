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
	var analyticsUrl = getAnalyticsUrlMap()
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   analyticsUrl[""],
		Path:   "api/v2/connections",
	}

	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	analyticdatalinks := AnalyticsDataDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &analyticdatalinks)

	return &analyticdatalinks, err
}

func (client *Client) CreateAnalyticsDataExport(ctx context.Context, analyticsdataToCreate AnalyticsDataCreateDto) (*AnalyticsDataDto, error) {
	var analyticsUrl = getAnalyticsUrlMap()
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   analyticsUrl[""],
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
	var analyticsUrl = getAnalyticsUrlMap()
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   analyticsUrl[""],
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
	var analyticsUrl = getAnalyticsUrlMap()
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   analyticsUrl[""],
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

func getAnalyticsUrlMap() map[string]string {
	return map[string]string{
		"US":   "https://na.csanalytics.powerplatform.microsoft.com/",
		"CAN":  "https://can.csanalytics.powerplatform.microsoft.com/",
		"SAM":  "https://sam.csanalytics.powerplatform.microsoft.com/",
		"EMEA": "https://emea.csanalytics.powerplatform.microsoft.com/",
		"OCE":  "https://oce.csanalytics.powerplatform.microsoft.com/",
		"PAC":  "https://apac.csanalytics.powerplatform.microsoft.com/",
		"JPN":  "https://jpn.csanalytics.powerplatform.microsoft.com/",
		"CHE":  "https://che.csanalytics.powerplatform.microsoft.com/",
		"FRA":  "https://fra.csanalytics.powerplatform.microsoft.com/",
		"UAE":  "https://uae.csanalytics.powerplatform.microsoft.com/",
		"GER":  "https://ger.csanalytics.powerplatform.microsoft.com/",
		"GBR":  "https://gbr.csanalytics.powerplatform.microsoft.com/",
		"IND":  "https://ind.csanalytics.powerplatform.microsoft.com/",
		"KOR":  "https://kor.csanalytics.powerplatform.microsoft.com/",
		"NOR":  "https://nor.csanalytics.powerplatform.microsoft.com/",
		"ZAF":  "https://zaf.csanalytics.powerplatform.microsoft.com/",
		"SGP":  "https://sgp.csanalytics.powerplatform.microsoft.com/",
		"SWE":  "https://swe.csanalytics.powerplatform.microsoft.com/",
		"GOV":  "https://gcc.csanalytics.powerplatform.microsoft.com/",
		"HIGH": "https://high.csanalytics.powerplatform.microsoft.com/",
	}
}
