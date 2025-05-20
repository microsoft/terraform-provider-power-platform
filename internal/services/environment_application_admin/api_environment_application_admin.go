// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_application_admin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newEnvironmentApplicationAdminClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) AddApplicationUser(ctx context.Context, environmentId, applicationId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/enroll"),
		RawQuery: url.Values{
			"api-version":   []string{"2020-10-01"},
			"environmentId": []string{environmentId},
			"appId":         []string{applicationId},
		}.Encode(),
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	return err
}

func (client *client) GetApplicationUser(ctx context.Context, environmentId, applicationId string) (bool, error) {
	// Build the Dataverse API endpoint URL
	// Use the PowerApps API to query the Dataverse environment
	dataverseApiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   fmt.Sprintf("%s.api.%s", environmentId, client.Api.GetConfig().Urls.PowerAppsUrl),
		Path:   fmt.Sprintf("/api/data/%s/applicationusers", constants.DATAVERSE_API_VERSION),
		RawQuery: url.Values{
			"$filter": []string{fmt.Sprintf("applicationid eq '%s'", applicationId)},
		}.Encode(),
	}

	var response applicationUserResponseDto
	apiResp, err := client.Api.Execute(
		ctx,
		[]string{client.Api.GetConfig().Urls.PowerAppsScope},
		"GET",
		dataverseApiUrl.String(),
		nil,
		nil,
		[]int{http.StatusOK},
		&response,
	)

	if err != nil {
		if apiResp != nil && apiResp.HttpResponse != nil && apiResp.HttpResponse.StatusCode == http.StatusNotFound {
			return false, nil // Environment or entity not found - not an error
		}
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			return false, nil
		}
		return false, err
	}

	// Check if we found the application user
	return len(response.Value) > 0, nil
}
