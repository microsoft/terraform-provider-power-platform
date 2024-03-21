// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewLanguagesClient(api *api.ApiClient) LanguagesClient {
	return LanguagesClient{
		Api: api,
	}
}

type LanguagesClient struct {
	Api *api.ApiClient
}

func (client *LanguagesClient) GetLanguagesByLocation(ctx context.Context, location string) (LanguagesDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentLanguages", location),
	}
	values := url.Values{
		"api-version": []string{"2023-06-01"},
	}
	apiUrl.RawQuery = values.Encode()

	languages := LanguagesDto{}

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return languages, err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &languages)

	if err != nil {
		return languages, err
	}

	return languages, nil
}
