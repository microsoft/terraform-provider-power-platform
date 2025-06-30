// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package analytics_data_export

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func NewAnalyticsExportClient(apiClient *api.Client, tenantClient tenant.Client) Client {
	return Client{
		Api:       apiClient,
		TenantApi: tenantClient,
	}
}

type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

// GetGatewayCluster retrieves information about a gateway cluster.
func (client *Client) GetGatewayCluster(ctx context.Context) (*GatewayClusterDto, error) {
	// Get tenant information using the tenant client
	tenantInfo, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant information: %w", err)
	}

	tenantApiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantInfo.TenantId, client.Api.Config.Urls.PowerPlatformUrl),
		Path:   "gateway/cluster",
		RawQuery: url.Values{
			"api-version": []string{"1"},
		}.Encode(),
	}

	gatewayCluster := GatewayClusterDto{}

	_, err = client.Api.Execute(ctx, nil, "GET", tenantApiUrl.String(), nil, nil, []int{http.StatusOK}, &gatewayCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetGatewayCluster API request: %w", err)
	}

	return &gatewayCluster, nil
}

func (client *Client) GetAnalyticsDataExport(ctx context.Context) ([]AnalyticsDataDto, error) {
	// Get the gateway cluster
	gatewayCluster, err := client.GetGatewayCluster(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gateway cluster: %w", err)
	}
	// Get the region from the gateway cluster
	region := gatewayCluster.GeoName

	// Get the analytics URL based on the region
	var analyticsUrl string
	analyticsUrl, err = getAnalyticsUrl(region)
	if err != nil {
		return nil, err
	}

	apiUrl, err := url.Parse(analyticsUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analytics URL: %w", err)
	}
	apiUrl.Path = "api/v2/connections"

	var adr AnalyticsDataResponse
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adr)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics data export: %w", err)
	}

	return adr.Value, nil
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
		"CH":   "https://che.csanalytics.powerplatform.microsoft.com/",
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
		"GOV":  "https://gcc.csanalytics.powerplatform.microsoft.us/",
		"HIGH": "https://high.csanalytics.powerplatform.microsoft.us/",
		"DOD":  "https://dod.csanalytics.csanalytics.appsplatform.us/",
	}
}

func getAnalyticsUrl(region string) (string, error) {
	urlMap := getAnalyticsUrlMap()
	analyticDataUrl, exists := urlMap[strings.ToUpper(region)]
	if !exists {
		return "", fmt.Errorf("invalid region: %s", region)
	}
	return analyticDataUrl, nil
}
