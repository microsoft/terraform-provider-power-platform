// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

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
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func NewEnvironmentClient(apiClient *api.Client) Client {
	return Client{
		tenantClient:   tenant.NewTenantClient(apiClient),
		solutionClient: solution.NewSolutionClient(apiClient),
		Api:            apiClient,
	}
}

type Client struct {
	tenantClient   tenant.Client
	solutionClient solution.Client
	Api            *api.Client
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

func (client *Client) GetLocations(ctx context.Context) (*LocationArrayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	locationsArray := LocationArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locationsArray)
	if err != nil {
		return nil, err
	}
	return &locationsArray, nil
}

func (client *Client) LocationValidator(ctx context.Context, location, azureRegion string) error {
	locationsArray, err := client.GetLocations(ctx)
	if err != nil {
		return err
	}

	foundLocation, err := findLocation(*locationsArray, location)
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

type currencyCodeValidatorDto struct {
	Name       string                             `json:"name"`
	ID         string                             `json:"id"`
	Type       string                             `json:"type"`
	Properties currencyCodeValidatorPropertiesDto `json:"properties"`
}

type currencyCodeValidatorPropertiesDto struct {
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}

type currencyCodeValidatorArrayDto struct {
	Value []currencyCodeValidatorDto `json:"value"`
}

func currencyCodeValidator(ctx context.Context, client *api.Client, location string, currencyCode string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	response, err := client.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	if err != nil {
		return err
	}

	defer response.HttpResponse.Body.Close()

	resp := currencyCodeValidatorArrayDto{}
	err = json.Unmarshal(response.BodyAsBytes, &resp)

	if err != nil {
		return err
	}

	codes := make([]string, len(resp.Value))
	for i, item := range resp.Value {
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

type languageCodeValidatorDto struct {
	Name       string                             `json:"name"`
	ID         string                             `json:"id"`
	Type       string                             `json:"type"`
	Properties languageCodeValidatorPropertiesDto `json:"properties"`
}

type languageCodeValidatorArrayDto struct {
	Value []languageCodeValidatorDto `json:"value"`
}

type languageCodeValidatorPropertiesDto struct {
	LocaleID        int    `json:"localeId"`
	LocalizedName   string `json:"localizedName"`
	DisplayName     string `json:"displayName"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}

func languageCodeValidator(ctx context.Context, client *api.Client, location string, languageCode string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentLanguages", location),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	response, err := client.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	if err != nil {
		return err
	}

	defer response.HttpResponse.Body.Close()

	resp := languageCodeValidatorArrayDto{}
	err = json.Unmarshal(response.BodyAsBytes, &resp)

	if err != nil {
		return err
	}

	codes := make([]string, len(resp.Value))
	for i, item := range resp.Value {
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

func (client *Client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	if environmentUrl == "" {
		return "", customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_ENVIRONMENT_URL_NOT_FOUND), "environment url not found, please check if the environment has dataverse linked")
	}

	envUrl, err := url.Parse(environmentUrl)
	if err != nil {
		return "", err
	}
	return envUrl.Host, nil
}

func (client *Client) GetEnvironment(ctx context.Context, environmentId string) (*EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("environment '%s' not found", environmentId))
		}
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata != nil && env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = constants.ZERO_UUID
	}

	if env.Properties.ParentEnvironmentGroup != nil && env.Properties.ParentEnvironmentGroup.Id == "" {
		env.Properties.ParentEnvironmentGroup.Id = constants.ZERO_UUID
	}

	return &env, nil
}

func (client *Client) DeleteEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	environmentDelete := enironmentDeleteDto{
		Code:    "7", // Application.
		Message: "Deleted using Power Platform Terraform Provider",
	}

	response, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, environmentDelete, []int{http.StatusNoContent, http.StatusAccepted, http.StatusConflict, http.StatusNotFound}, nil)

	// Handle HTTP 404 case - if the environment is not found, consider it already deleted
	if response != nil && response.HttpResponse.StatusCode == http.StatusNotFound {
		tflog.Info(ctx, fmt.Sprintf("Environment '%s' not found. Treating as successfully deleted.", environmentId))
		return nil
	}

	if response.HttpResponse.StatusCode == http.StatusConflict {
		err := client.handleHttpConflict(ctx, response)
		if err != nil {
			return err
		}
		return client.DeleteEnvironment(ctx, environmentId)
	}

	var httpError *customerrors.UnexpectedHttpStatusCodeError
	if errors.As(err, &httpError) {
		return fmt.Errorf("unexpected HTTP Status %s; Body: %s", httpError.StatusText, httpError.Body)
	}

	tflog.Debug(ctx, "Environment Deletion Operation HTTP Status: '"+response.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for environment deletion operation to complete")

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, response)
	if err != nil {
		return err
	}

	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Environment deletion failed. Retrying")
		return client.DeleteEnvironment(ctx, environmentId)
	}
	return nil
}

func (client *Client) AddDataverseToEnvironment(ctx context.Context, environmentId string, environmentCreateLinkEnvironmentMetadata createLinkEnvironmentMetadataDto) (*EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/provisionInstance", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentCreateLinkEnvironmentMetadata, []int{http.StatusAccepted}, nil)
	if err != nil {
		tflog.Error(ctx, "Error adding Dataverse to environment: "+err.Error())
	}

	tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")

	locationHeader := apiResponse.GetHeader(constants.HEADER_LOCATION)
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	_, err = url.Parse(locationHeader)
	if err != nil {
		tflog.Error(ctx, "Error parsing location header: "+err.Error())
	}

	retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
	tflog.Debug(ctx, "Retry Header: "+retryHeader)
	retryAfter, err := time.ParseDuration(retryHeader)
	if err != nil {
		retryAfter = api.DefaultRetryAfter()
	} else {
		retryAfter = retryAfter * time.Second
	}
	for {
		lifecycleEnv := &EnvironmentDto{}
		lifecycleResponse, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted, http.StatusConflict}, &lifecycleEnv)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation Operation HTTP Status: '%s'", lifecycleResponse.HttpResponse.Status))
		if lifecycleResponse.HttpResponse.StatusCode == http.StatusConflict {
			continue
		}

		if lifecycleEnv == nil || lifecycleEnv.Properties == nil {
			tflog.Debug(ctx, fmt.Sprintf("The environment lifecycle response body did not match expected format. Response status code: %s", lifecycleResponse.HttpResponse.Status))
			continue
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Dataverse Creation Operation State: '%s'", lifecycleEnv.Properties.ProvisioningState))

		if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
			return lifecycleEnv, nil
		} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
			return lifecycleEnv, fmt.Errorf("dataverse creation failed. provisioning state: %s", lifecycleEnv.Properties.ProvisioningState)
		}
	}
}

func (client *Client) ModifyEnvironmentType(ctx context.Context, environmentId, environmentType string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/modifySku", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	modifySkuDto := modifySkuDto{
		EnvironmentSku: environmentType,
	}

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, modifySkuDto, []int{http.StatusAccepted, http.StatusOK, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}

	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Environment update failed. Retrying")
		return client.ModifyEnvironmentType(ctx, environmentId, environmentType)
	}
	return nil
}

func (client *Client) CreateEnvironment(ctx context.Context, environmentToCreate environmentCreateDto) (*EnvironmentDto, error) {
	if environmentToCreate.Properties.LinkedEnvironmentMetadata != nil && environmentToCreate.Location != "" && environmentToCreate.Properties.LinkedEnvironmentMetadata.DomainName != "" {
		err := client.ValidateCreateEnvironmentDetails(ctx, environmentToCreate.Location, environmentToCreate.Properties.LinkedEnvironmentMetadata.DomainName)
		if err != nil {
			return nil, err
		}
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentToCreate, []int{http.StatusAccepted, http.StatusCreated, http.StatusInternalServerError, http.StatusConflict}, nil)
	if err != nil {
		return nil, err
	}

	if apiResponse.HttpResponse.StatusCode == http.StatusConflict {
		err := client.handleHttpConflict(ctx, apiResponse)
		if err != nil {
			return nil, err
		}
		return client.CreateEnvironment(ctx, environmentToCreate)
	}

	if apiResponse.HttpResponse.StatusCode == http.StatusInternalServerError {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_ENVIRONMENT_CREATION), string(apiResponse.BodyAsBytes))
	}

	tflog.Debug(ctx, "Environment Creation Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")

	createdEnvironmentId := ""
	switch apiResponse.HttpResponse.StatusCode {
	case http.StatusAccepted:
		lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
		if err != nil {
			return nil, err
		}

		if lifecycleResponse.State.Id == "Succeeded" {
			parts := strings.Split(lifecycleResponse.Links.Environment.Path, "/")
			if len(parts) == 0 {
				return nil, errors.New("can't parse environment id from response " + lifecycleResponse.Links.Environment.Path)
			}
			createdEnvironmentId = parts[len(parts)-1]
			tflog.Debug(ctx, "Created Environment Id: "+createdEnvironmentId)
		}

	case http.StatusCreated:
		envCreatedResponse := lifecycleCreatedDto{}
		err := apiResponse.MarshallTo(&envCreatedResponse)
		if err != nil {
			return nil, err
		}
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

func (client *Client) UpdateEnvironmentAiFeatures(ctx context.Context, environmentId string, generativeAIConfig GenerativeAiFeaturesDto) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, generativeAIConfig, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	if apiResponse.HttpResponse.StatusCode == http.StatusConflict {
		err := client.handleHttpConflict(ctx, apiResponse)
		if err != nil {
			return err
		}
		return client.UpdateEnvironmentAiFeatures(ctx, environmentId, generativeAIConfig)
	}

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}

	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Environment update ai features failed. Retrying")
		return client.UpdateEnvironmentAiFeatures(ctx, environmentId, generativeAIConfig)
	}
	return nil
}

func (client *Client) handleHttpConflict(ctx context.Context, apiResponse *api.Response) error {
	body := string(apiResponse.BodyAsBytes)
	if body == "" {
		return errors.New("environment failed with HTTP 409. No body in response")
	}
	// if 409 returns anything other than another ongoing lifecycle operation, fail the request and return the body as error to the user
	if !strings.Contains(body, "OperationNotStartable") {
		return errors.New("environment failed with HTTP 409. Body: " + body)
	}
	tflog.Debug(ctx, "Another lifecycle operation is in progress, waiting for it to complete")
	return client.Api.SleepWithContext(ctx, api.DefaultRetryAfter())
}

func (client *Client) UpdateEnvironment(ctx context.Context, environmentId string, environment EnvironmentDto) (*EnvironmentDto, error) {
	if environment.Location != "" && environment.Properties.LinkedEnvironmentMetadata != nil && environment.Properties.LinkedEnvironmentMetadata.DomainName != "" {
		err := client.ValidateUpdateEnvironmentDetails(ctx, environment.Id, environment.Properties.LinkedEnvironmentMetadata.DomainName)
		if err != nil {
			return nil, err
		}
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	// Due to a bug in BAPI that triggers managed environment on update of a description field, we need to use the older API version
	// values.Add("api-version", "2022-05-01")
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()
	apiResponse, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, environment, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return nil, err
	}

	if apiResponse.HttpResponse.StatusCode == http.StatusConflict {
		err := client.handleHttpConflict(ctx, apiResponse)
		if err != nil {
			return nil, err
		}
		return client.UpdateEnvironment(ctx, environmentId, environment)
	}

	// wait for the lifecycle operation to finish.
	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return nil, err
	}

	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
		tflog.Info(ctx, "Environment update failed. Retrying")
		return client.UpdateEnvironment(ctx, environmentId, environment)
	}

	// despite lifecycle operation success, the environment may not be ready yet.
	for {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
		env, err := client.GetEnvironment(ctx, environmentId)
		if err != nil {
			return nil, err
		}
		tflog.Info(ctx, "Environment State: '"+env.Properties.States.Management.Id+"'")
		if env.Properties.States.Management.Id == "Ready" {
			return env, nil
		} else if env.Properties.States.Management.Id == "Running" {
			continue
		}
		return nil, errors.New("environment update failed. unexpected management state: " + env.Properties.States.Management.Id)
	}
}

func (client *Client) GetEnvironments(ctx context.Context) ([]EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("$expand", "properties/billingPolicy,properties/copilotPolicies")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	envArray := environmentArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}

func (client *Client) GetDefaultCurrencyForEnvironment(ctx context.Context, environmentId string) (*TransactionCurrencyDto, error) {
	orgSettings := organizationSettingsArrayDto{}
	err := client.solutionClient.GetTableData(ctx, environmentId, "organizations", "", &orgSettings)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)

	currencies := transactionCurrencyArrayDto{}
	err = client.solutionClient.GetTableData(ctx, environmentId, "transactioncurrencies", values.Encode(), &currencies)
	if err != nil {
		return nil, err
	}
	if len(currencies.Value) == 0 {
		return nil, fmt.Errorf("no default currency found for environment %s", environmentId)
	}
	return &currencies.Value[0], nil
}

func (client *Client) ValidateCreateEnvironmentDetails(ctx context.Context, location, domain string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails",
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	envDetails := validateCreateEnvironmentDetailsDto{
		DomainName:          domain,
		EnvironmentLocation: location,
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, envDetails, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ValidateUpdateEnvironmentDetails(ctx context.Context, environmentId, domain string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails",
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	envDetails := validateUpdateEnvironmentDetailsDto{
		DomainName:      domain,
		EnvironmentName: environmentId,
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, envDetails, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}
