// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_wave

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

type environmentWaveClient struct {
	Api *api.Client
}

func newEnvironmentWaveClient(apiClient *api.Client) *environmentWaveClient {
	return &environmentWaveClient{
		Api: apiClient,
	}
}

func (client *environmentWaveClient) GetGeoFromEnvironment(ctx context.Context, environmentId string) (*string, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.AdminPowerPlatformUrl,
		Path:   "/api/tenants/mytenant/organizations",
	}

	organizations := OrganizationsArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &organizations)
	if err != nil {
		return nil, err
	}

	for _, org := range organizations {
		if org.Id == environmentId {
			return &org.CrmGeo, nil
		}
	}
	return nil, fmt.Errorf("geo for environment with ID %s not found", environmentId)
}

func (client *environmentWaveClient) UpdateFeature(ctx context.Context, environmentId string, featureName string) (*FeatureDto, error) {
	geo, err := client.GetGeoFromEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.AdminPowerPlatformUrl,
		Path:   fmt.Sprintf("/api/environments/%s/features/%s/enable", environmentId, featureName),
	}

	values := url.Values{}
	values.Add("geo", *geo)
	apiUrl.RawQuery = values.Encode()

	urlString := apiUrl.String()

	_, err = client.Api.Execute(ctx, nil, "POST", urlString, nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	retryAfter := api.DefaultRetryAfter()
	for {
		feature, err := client.GetFeature(ctx, environmentId, featureName)
		if err != nil {
			return nil, err
		}

		if feature != nil && feature.AppsUpgradeState != "Upgrading" {
			tflog.Info(ctx, fmt.Sprintf("Feature %s  with state: %s", featureName, feature.AppsUpgradeState))
			return feature, nil
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Feature %s not yet enabled, polling...", featureName))
	}
}

func (client *environmentWaveClient) GetFeature(ctx context.Context, environmentId string, featureName string) (*FeatureDto, error) {
	geo, err := client.GetGeoFromEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.AdminPowerPlatformUrl,
		Path:   fmt.Sprintf("/api/environments/%s/features", environmentId),
	}

	values := url.Values{}
	values.Add("geo", *geo)
	apiUrl.RawQuery = values.Encode()

	features := FeaturesArrayDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &features)
	if err != nil {
		return nil, err
	}

	for _, feature := range features.Values {
		if feature.FeatureName == featureName {
			return &feature, nil
		}
	}

	return nil, nil
}
