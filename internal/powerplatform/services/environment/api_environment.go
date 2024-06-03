// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
)

func NewEnvironmentClient(api *api.ApiClient) EnvironmentClient {
	return EnvironmentClient{
		solutionClient: solution.NewSolutionClient(api),
		Api:            api,
	}
}

type EnvironmentClient struct {
	solutionClient solution.SolutionClient
	Api            *api.ApiClient
}

func findLocation(locations LocationArrayDto, locationToFind string) (*LocationDto, error) {
	for _, loc := range locations.Value {
		if loc.Name == locationToFind {
			return &loc, nil
		}
	}

	locationNames := make([]string, len(locations.Value))
	for i, loc := range locations.Value {
		locationNames[i] = loc.Name
	}
	return nil, fmt.Errorf("location '%s' is not valid. valid locations are: %s", locationToFind, strings.Join(locationNames, ", "))
}

func findAzureRegion(location *LocationDto, azureRegion string) (bool, error) {
	for _, region := range location.Properties.AzureRegions {
		if region == azureRegion {
			return true, nil
		}
	}
	return false, fmt.Errorf("region '%s' is not valid for location %s. valid regions are: %s", azureRegion, location.Name, strings.Join(location.Properties.AzureRegions, ", "))
}

func locationValidator(client *api.ApiClient, location, azureRegion string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	response, err := client.Execute(context.Background(), "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	if err != nil {
		return err
	}

	defer response.Response.Body.Close()

	locationsArray := LocationArrayDto{}
	err = json.Unmarshal(response.BodyAsBytes, &locationsArray)

	if err != nil {
		return err
	}

	foundLocation, err := findLocation(locationsArray, location)
	if err != nil {
		return err
	}

	if azureRegion == "" {
		return nil
	}

	isRegionFound, err := findAzureRegion(foundLocation, azureRegion)
	if err != nil || !isRegionFound {
		return err
	}

	return nil
}

func currencyCodeValidator(client *api.ApiClient, location string, currencyCode string) error {
	var parsed struct {
		Value []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Type       string `json:"type"`
			Properties struct {
				Code            string `json:"code"`
				Symbol          string `json:"symbol"`
				IsTenantDefault bool   `json:"isTenantDefault"`
			} `json:"properties"`
		} `json:"value"`
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	response, err := client.Execute(context.Background(), "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	if err != nil {
		return err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &parsed)

	if err != nil {
		return err
	}

	codes := make([]string, len(parsed.Value))
	for i, item := range parsed.Value {
		codes[i] = item.Name
	}

	found := func(items []string, check string) bool {
		for _, item := range items {
			if item == check {
				return true
			}
		}
		return false
	}(codes, currencyCode)

	if !found {
		return fmt.Errorf("currency Code %s is not valid. valid currency codes are: %s", currencyCode, strings.Join(codes, ", "))
	}

	return nil
}

func languageCodeValidator(client *api.ApiClient, location string, languageCode string) error {
	var parsed struct {
		Value []struct {
			Name       string `json:"name"`
			ID         string `json:"id"`
			Type       string `json:"type"`
			Properties struct {
				LocaleID        int    `json:"localeId"`
				LocalizedName   string `json:"localizedName"`
				DisplayName     string `json:"displayName"`
				IsTenantDefault bool   `json:"isTenantDefault"`
			} `json:"properties"`
		} `json:"value"`
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentLanguages", location),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	response, err := client.Execute(context.Background(), "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	if err != nil {
		return err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &parsed)

	if err != nil {
		return err
	}

	codes := make([]string, len(parsed.Value))
	for i, item := range parsed.Value {
		codes[i] = item.Name
	}

	found := func(items []string, check string) bool {
		for _, item := range items {
			if item == check {
				return true
			}
		}
		return false
	}(codes, languageCode)

	if !found {
		return fmt.Errorf("language Code %s is not valid. valid language codes are: %s", languageCode, strings.Join(codes, ", "))
	}

	return nil
}

func (client *EnvironmentClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *EnvironmentClient) GetEnvironment(ctx context.Context, environmentId string) (*EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, powerplatform_helpers.WrapIntoProviderError(err, powerplatform_helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("environment '%s' not found", environmentId))
		}
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata != nil && env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *EnvironmentClient) DeleteEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	environmentDelete := EnvironmentDeleteDto{
		Code:    "7", //Application
		Message: "Deleted using Power Platform Terraform Provider",
	}

	response, err := client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, environmentDelete, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}
	tflog.Debug(ctx, "Environment Deletion Operation HTTP Status: '"+response.Response.Status+"'")

	tflog.Debug(ctx, "Waiting for environment deletion operation to complete")
	_, err = client.Api.DoWaitForLifecycleOperationStatus(ctx, response)
	if err != nil {
		return err
	}
	return nil
}

func (client *EnvironmentClient) AddDataverseToEnvironment(ctx context.Context, environmentId string, environmentCreateLinkEnvironmentMetadata EnvironmentCreateLinkEnvironmentMetadataDto) (*EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/provisionInstance", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	apiResponse, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environmentCreateLinkEnvironmentMetadata, []int{http.StatusAccepted}, nil)
	if err != nil {
		tflog.Error(ctx, "Error adding Dataverse to environment: "+err.Error())
	}

	tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.Response.Status+"'")

	locationHeader := apiResponse.GetHeader("Location")
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	_, err = url.Parse(locationHeader)
	if err != nil {
		tflog.Error(ctx, "Error parsing location header: "+err.Error())
	}

	retryHeader := apiResponse.GetHeader("Retry-After")
	tflog.Debug(ctx, "Retry Header: "+retryHeader)
	retryAfter, err := time.ParseDuration(retryHeader)
	if err != nil {
		retryAfter = time.Duration(5) * time.Second
	} else {
		retryAfter = retryAfter * time.Second
	}
	for {
		lifecycleEnv := EnvironmentDto{}
		lifecycleResponse, err := client.Api.Execute(ctx, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, &lifecycleEnv)
		if err != nil {
			return nil, err
		}

		client.Api.Sleep(retryAfter)

		tflog.Debug(ctx, "Dataverse Creation Operation State: '"+lifecycleEnv.Properties.ProvisioningState+"'")
		tflog.Debug(ctx, "Dataverse Creation Operation HTTP Status: '"+lifecycleResponse.Response.Status+"'")

		if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
			return &lifecycleEnv, nil
		} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
			return &lifecycleEnv, errors.New("dataverse creation failed. provisioning state: " + lifecycleEnv.Properties.ProvisioningState)
		}
	}
}

func (client *EnvironmentClient) CreateEnvironment(ctx context.Context, environmentToCreate EnvironmentCreateDto) (*EnvironmentDto, error) {
	if environmentToCreate.Properties.LinkedEnvironmentMetadata != nil && environmentToCreate.Location != "" && environmentToCreate.Properties.LinkedEnvironmentMetadata.DomainName != "" {
		err := client.ValidateEnvironmentDetails(ctx, environmentToCreate.Location, environmentToCreate.Properties.LinkedEnvironmentMetadata.DomainName)
		if err != nil {
			return nil, err
		}
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environmentToCreate, []int{http.StatusAccepted, http.StatusCreated}, nil)
	if err != nil {
		return nil, err
	}

	tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.Response.Status+"'")

	createdEnvironmentId := ""
	if apiResponse.Response.StatusCode == http.StatusAccepted {

		lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
		if err != nil {
			return nil, err
		}

		if lifecycleResponse.State.Id == "Succeeded" {
			parts := strings.Split(lifecycleResponse.Links.Environment.Path, "/")
			if len(parts) > 0 {
				createdEnvironmentId = parts[len(parts)-1]
			} else {
				return nil, errors.New("can't parse environment id from response " + lifecycleResponse.Links.Environment.Path)
			}
			tflog.Debug(ctx, "Created Environment Id: "+createdEnvironmentId)
		}
	} else if apiResponse.Response.StatusCode == http.StatusCreated {
		envCreatedResponse := EnvironmentLifecycleCreatedDto{}
		apiResponse.MarshallTo(&envCreatedResponse)
		if envCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return nil, errors.New("environment creation failed. provisioning state: " + envCreatedResponse.Properties.ProvisioningState)
		}
		createdEnvironmentId = envCreatedResponse.Name
	}

	env, err := client.GetEnvironment(ctx, createdEnvironmentId)
	if err != nil {
		return &EnvironmentDto{}, fmt.Errorf("environment '%s' not found. '%s'", createdEnvironmentId, err)
	}
	if env.Properties.LinkedEnvironmentMetadata != nil && environmentToCreate.Properties.LinkedEnvironmentMetadata != nil && environmentToCreate.Properties.LinkedEnvironmentMetadata.Templates != nil {
		env.Properties.LinkedEnvironmentMetadata.Templates = environmentToCreate.Properties.LinkedEnvironmentMetadata.Templates
		env.Properties.LinkedEnvironmentMetadata.TemplateMetadata = environmentToCreate.Properties.LinkedEnvironmentMetadata.TemplateMetadata
	}
	return env, err
}

func (client *EnvironmentClient) UpdateEnvironment(ctx context.Context, environmentId string, environment EnvironmentDto) (*EnvironmentDto, error) {
	if environment.Location != "" && environment.Properties.LinkedEnvironmentMetadata != nil && environment.Properties.LinkedEnvironmentMetadata.DomainName != "" {
		err := client.ValidateEnvironmentDetails(ctx, environment.Location, environment.Properties.LinkedEnvironmentMetadata.DomainName)
		if err != nil {
			return nil, err
		}
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2022-05-01")
	apiUrl.RawQuery = values.Encode()
	_, err := client.Api.Execute(ctx, "PATCH", apiUrl.String(), nil, environment, []int{http.StatusAccepted}, nil)
	if err != nil {
		return nil, err
	}

	client.Api.Sleep(10 * time.Second)

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
				client.Api.Sleep(3 * time.Second)
				if createdEnv.Properties.States.Management.Id == "Ready" {

					return createdEnv, nil
				}

			}
		}
	}

	return nil, fmt.Errorf("environment '%s' not found", environmentId)
}

func (client *EnvironmentClient) GetEnvironments(ctx context.Context) ([]EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("$expand", "properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	envArray := EnvironmentDtoArray{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *EnvironmentClient) GetDefaultCurrencyForEnvironment(ctx context.Context, environmentId string) (*TransactionCurrencyDto, error) {
	orgSettings := OrganizationSettingsArrayDto{}
	err := client.solutionClient.GetTableData(ctx, environmentId, "organizations", "", &orgSettings)
	if err != nil {
		return nil, err
	} else {
		values := url.Values{}
		values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)

		currencies := TransactionCurrencyArrayDto{}
		err := client.solutionClient.GetTableData(ctx, environmentId, "transactioncurrencies", values.Encode(), &currencies)
		if err != nil {
			return nil, err
		} else {
			if currencies.Value != nil && len(currencies.Value) >= 1 {
				return &currencies.Value[0], nil
			} else {
				return nil, fmt.Errorf("no default currency found for environment %s", environmentId)
			}
		}
	}
}

func (client *EnvironmentClient) ValidateEnvironmentDetails(ctx context.Context, location, domain string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails",
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	envDetails := ValidateEnvironmentDetailsDto{
		DomainName:          domain,
		EnvironmentLocation: location,
	}

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, envDetails, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}
