// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestUnitDoWaitForLifecycleOperationStatus_NoLocationHeader(t *testing.T) {
	t.Parallel()

	client := api.NewApiClientBase(&config.ProviderConfig{}, api.NewAuthBase(&config.ProviderConfig{}))
	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), resp)
	require.NoError(t, err)
	require.Nil(t, lifecycle)
}

func TestUnitDoWaitForLifecycleOperationStatus_InvalidLocationHeader(t *testing.T) {
	t.Parallel()

	clientCfg := config.ProviderConfig{Urls: config.ProviderConfigUrls{
		BapiUrl:            "https://example.invalid",
		PowerAppsUrl:       "https://example.invalid",
		PowerPlatformUrl:   "https://example.invalid",
		PowerAppsScope:     "scope",
		PowerPlatformScope: "scope",
	}}
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&clientCfg))

	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{clientCfg.Urls.BapiUrl + "/%zz"},
		constants.HEADER_RETRY_AFTER: []string{"0"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), resp)
	require.Error(t, err)
	var urlErr *customerrors.UrlFormatError
	require.ErrorAs(t, err, &urlErr)
	require.Nil(t, lifecycle)
}

// lifecycleTestLocation is an https URL on the BAPI host, because a request without explicit scopes is only
// authorized for https URLs of known hosts. httpmock serves it, so the lifecycle tests that poll it cannot run
// in parallel with the tests that use the real http.DefaultTransport.
const lifecycleTestLocation = "https://" + constants.PUBLIC_BAPI_DOMAIN + "/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001"

func lifecycleTestUrls() config.ProviderConfigUrls {
	return config.ProviderConfigUrls{
		BapiUrl:        constants.PUBLIC_BAPI_DOMAIN,
		PowerAppsScope: constants.PUBLIC_POWERAPPS_SCOPE,
	}
}

func TestUnitDoWaitForLifecycleOperationStatus_ConflictThenSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	callCount := 0
	httpmock.RegisterResponder(http.MethodGet, lifecycleTestLocation, func(*http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return httpmock.NewStringResponse(http.StatusConflict, ""), nil
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"state":{"id":"Succeeded"}}`), nil
	})

	clientCfg := config.ProviderConfig{
		TestMode: true,
		Urls:     lifecycleTestUrls(),
	}
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&clientCfg))

	ctx := helpers.TestContext(context.Background(), t.Name())
	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_OPERATION_LOCATION: []string{lifecycleTestLocation},
		constants.HEADER_RETRY_AFTER:        []string{"0"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_SucceededNoSleep(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, lifecycleTestLocation,
		httpmock.NewStringResponder(http.StatusOK, `{"state":{"id":"Succeeded"}}`))

	clientCfg := config.ProviderConfig{
		TestMode: false,
		Urls:     lifecycleTestUrls(),
	}
	authCfg := clientCfg
	authCfg.TestMode = true
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{lifecycleTestLocation},
		constants.HEADER_RETRY_AFTER: []string{"120"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextTimeout(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, lifecycleTestLocation,
		httpmock.NewStringResponder(http.StatusOK, `{"state":{"id":"InProgress"}}`))

	clientCfg := config.ProviderConfig{
		TestMode: false,
		Urls:     lifecycleTestUrls(),
	}
	authCfg := clientCfg
	// Use TestMode=false for lifecycle polling, but keep auth in TestMode to avoid real token requests.
	authCfg.TestMode = true
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{lifecycleTestLocation},
		constants.HEADER_RETRY_AFTER: []string{"120"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, lifecycle)
}
