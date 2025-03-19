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

func (client *environmentWaveClient) UpdateFeature(ctx context.Context, environmentId string, featureName string) (*FeatureDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   "api.admin.powerplatform.microsoft.com",
		Path:   fmt.Sprintf("/api/environments/%s/features/%s/enable", environmentId, featureName),
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	retryAfter := api.DefaultRetryAfter()
	for {
		feature, err := client.GetFeature(ctx, environmentId, featureName)
		if err != nil {
			return nil, err
		}

		if feature.Enabled {
			tflog.Info(ctx, fmt.Sprintf("Feature %s enabled with state: %s", featureName, feature.AppsUpgradeState))
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
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   "api.admin.powerplatform.microsoft.com",
		Path:   fmt.Sprintf("/api/environments/%s/features", environmentId),
	}

	features := FeaturesArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &features)
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
