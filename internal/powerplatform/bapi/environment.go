package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2022-05-01", nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	envArray := models.EnvironmentDtoArray{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *ApiClient) GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+environmentId+"?$expand=permissions,properties.capacity&api-version=2022-05-01", nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	env := models.EnvironmentDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if err != nil {
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *ApiClient) DeleteEnvironment(ctx context.Context, environmentId string) error {
	request, err := http.NewRequestWithContext(ctx, "DELETE", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+environmentId+"?api-version=2020-06-01", nil)

	if err != nil {
		return err
	}

	_, err = client.doRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *ApiClient) CreateEnvironment(ctx context.Context, environment models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
	body, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2022-05-01", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	_, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	environments, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	for _, env := range environments {
		if env.Location == environment.Location && env.Properties.DisplayName == environment.Properties.DisplayName {
			for {
				createdEnv, err := client.GetEnvironment(ctx, env.Name)
				if err != nil {
					return nil, err
				}
				tflog.Info(ctx, "Environment State: '"+createdEnv.Properties.States.Management.Id+"'")
				time.Sleep(1 * time.Second)
				if createdEnv.Properties.States.Management.Id != "Running" {
					return createdEnv, nil
				}

			}
		}
	}
	return &models.EnvironmentDto{}, errors.New("environment not found")
}

func (client *ApiClient) UpdateEnvironment(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error) {
	body, err := json.Marshal(environment)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+environmentId+"?api-version=2022-05-01", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	_, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	environments, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	for _, env := range environments {
		if env.Name == environmentId {
			for {
				createdEnv, err := client.GetEnvironment(ctx, env.Name)
				if err != nil {
					return nil, err
				}
				tflog.Info(ctx, "Environment State: '"+createdEnv.Properties.States.Management.Id+"'")
				time.Sleep(3 * time.Second)
				if createdEnv.Properties.States.Management.Id == "Ready" {

					return createdEnv, nil
				}

			}
		}
	}

	return nil, errors.New("Environment not found")
}
