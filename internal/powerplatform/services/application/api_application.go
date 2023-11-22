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

	application := ApplicationArrayDto{}

	_, err := client.baseApi.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *ApplicationClient) InstallApplicationInEnvironment(ctx context.Context, environmentId string, uniqueName string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.baseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages/%s/install", environmentId, uniqueName),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	_, err := client.baseApi.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	return nil
}
