package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewApplicationClient(powerPlatformClientApi *api.PowerPlatformClientApi) ApplicationClient {
	return ApplicationClient{
		baseApi: powerPlatformClientApi,
	}
}

type ApplicationClient struct {
	baseApi *api.PowerPlatformClientApi
}

func (client *ApplicationClient) GetApplicationsByEnvironmentId(ctx context.Context, environmentId string) ([]ApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.baseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages", environmentId),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	header := http.Header{
		"$expand": []string{"permissions,properties.capacity"},
	}

	method := "GET"

	application := []ApplicationDto{}

	_, err := client.baseApi.Execute(ctx, method, apiUrl.String(), header, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application, nil
}
