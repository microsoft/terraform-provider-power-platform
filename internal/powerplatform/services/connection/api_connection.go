// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewConnectionsClient(api *api.ApiClient) ConnectionsClient {
	return ConnectionsClient{
		Api: api,
	}
}

type ConnectionsClient struct {
	Api *api.ApiClient
}

func (client *ConnectionsClient) BuildHostUri(environmentId string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, client.Api.GetConfig().Urls.PowerPlatformUrl)

}

func (client *ConnectionsClient) GetConnections(ctx context.Context, environmentId string) ([]ConnectionDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   "/connectivity/connections",
	}

	values := url.Values{}
	//values.Add("$expand", "permissions($filter=maxAssignedTo('<<aaid_objectid_here>>'))")
	//values.Add("$filter", fmt.Sprintf("environment eq '%s' and ApiId not in ('shared_logicflows', 'shared_powerflows', 'shared_pqogenericconnector')", environmentId))
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	connetionsArray := ConnectionDtoArray{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{200}, &connetionsArray)
	if err != nil {
		return nil, err
	}

	return connetionsArray.Value, nil
}
