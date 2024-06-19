// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	u "net/url"
	"strings"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewWebApiClient(api *api.ApiClient) WebApiClient {
	return WebApiClient{
		Api: api,
	}
}

type WebApiClient struct {
	Api *api.ApiClient
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}

func (client *WebApiClient) ExecuteWebApiRequest(ctx context.Context, environmentId, url, method string, body *string, headers map[string]string) (*api.ApiHttpResponse, error) {
	environmentUrl, err := client.getEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	e, _ := u.Parse(environmentUrl)
	apiUrl := fmt.Sprintf("%s://%s/%s", e.Scheme, e.Host, url)

	h := http.Header{}
	for k, v := range headers {
		h.Add(k, v)
	}

	res, err := client.Api.Execute(ctx, method, apiUrl, h, body, nil, nil)
	return res, err
}

func (client *WebApiClient) getEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *WebApiClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
