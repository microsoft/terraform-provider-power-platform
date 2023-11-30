package powerplatform

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (client *ApplicationClient) InstallApplicationInEnvironment(ctx context.Context, environmentId string, uniqueName string) (string, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.baseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages/%s/install", environmentId, uniqueName),
	}
	values := url.Values{
		"api-version": []string{"2022-03-01-preview"},
	}
	apiUrl.RawQuery = values.Encode()

	response, err := client.baseApi.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return "", err
	}

	applicationId := ""
	if response.Response.StatusCode == http.StatusAccepted {
		locationHeader := response.GetHeader("Location")
		tflog.Debug(ctx, "Location Header: "+locationHeader)

		_, err = url.Parse(locationHeader)
		if err != nil {
			tflog.Error(ctx, "Error parsing location header: "+err.Error())
		}

		retryHeader := response.GetHeader("Retry-After")
		tflog.Debug(ctx, "Retry Header: "+retryHeader)

		retryAfter, err := time.ParseDuration(retryHeader)
		if err != nil {
			retryAfter = time.Duration(5) * time.Second
		} else {
			retryAfter = retryAfter * time.Second
		}

		for {
			lifecycleResponse := ApplicationLifecycleDto{}
			_, err = client.baseApi.Execute(ctx, "GET", locationHeader, nil, nil, []int{http.StatusOK}, &lifecycleResponse)
			if err != nil {
				return "", err
			}

			time.Sleep(retryAfter)

			if lifecycleResponse.State.Id == "Succeeded" {
				parts := strings.Split(lifecycleResponse.Links.Environment.Path, "/")
				if len(parts) > 0 {
					applicationId = parts[len(parts)-1]
				} else {
					return "", errors.New("can't parse environment id from response " + lifecycleResponse.Links.Environment.Path)
				}
				tflog.Debug(ctx, "Created Environment Id: "+applicationId)
				break
			}
		}
	} else if response.Response.StatusCode == http.StatusCreated {
		envCreatedResponse := ApplicationLifecycleCreatedDto{}
		response.MarshallTo(&envCreatedResponse)
		if envCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return "", errors.New("environment creation failed. provisioning state: " + envCreatedResponse.Properties.ProvisioningState)
		}
		applicationId = envCreatedResponse.Name
	}

	return applicationId, nil
}
