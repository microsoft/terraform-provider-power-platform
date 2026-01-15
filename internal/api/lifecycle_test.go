// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestUnitDoWaitForLifecycleOperationStatus_NoLocationHeader(t *testing.T) {
	t.Parallel()

	client := NewApiClientBase(&config.ProviderConfig{}, nil)
	resp := &Response{HttpResponse: &http.Response{Header: http.Header{}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), resp)
	require.NoError(t, err)
	require.Nil(t, lifecycle)
}

func TestUnitDoWaitForLifecycleOperationStatus_ParseErrorAndExecuteError(t *testing.T) {
	t.Parallel()

	clientCfg := config.ProviderConfig{Urls: config.ProviderConfigUrls{
		BapiUrl:            "https://example.invalid",
		PowerAppsUrl:       "https://example.invalid",
		PowerPlatformUrl:   "https://example.invalid",
		PowerAppsScope:     "scope",
		PowerPlatformScope: "scope",
	}}
	client := NewApiClientBase(&clientCfg, nil)

	resp := &Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{clientCfg.Urls.BapiUrl + "/%zz"},
		constants.HEADER_RETRY_AFTER: []string{"0"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), resp)
	require.Error(t, err)
	require.Nil(t, lifecycle)
}

func TestUnitDoWaitForLifecycleOperationStatus_ConflictThenSuccess(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"state":{"id":"Succeeded"}}`))
	}))
	defer server.Close()

	clientCfg := config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl:            server.URL,
			PowerAppsUrl:       server.URL,
			PowerPlatformUrl:   server.URL,
			PowerAppsScope:     "scope",
			PowerPlatformScope: "scope",
		},
	}
	client := NewApiClientBase(&clientCfg, NewAuthBase(&clientCfg))

	ctx := helpers.TestContext(context.Background(), t.Name())
	location := server.URL + "/lifecycle/status"
	resp := &Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_OPERATION_LOCATION: []string{location},
		constants.HEADER_RETRY_AFTER:        []string{"0"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_SleepWithContextError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"state":{"id":"InProgress"}}`))
	}))
	defer server.Close()

	clientCfg := config.ProviderConfig{
		TestMode: false,
		Urls: config.ProviderConfigUrls{
			BapiUrl:            server.URL,
			PowerAppsUrl:       server.URL,
			PowerPlatformUrl:   server.URL,
			PowerAppsScope:     "scope",
			PowerPlatformScope: "scope",
		},
	}
	authCfg := clientCfg
	authCfg.TestMode = true
	client := NewApiClientBase(&clientCfg, NewAuthBase(&authCfg))

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Millisecond, cancel)

	resp := &Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{server.URL + "/lifecycle/status"},
		constants.HEADER_RETRY_AFTER: []string{"120"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.Error(t, err)
	require.Nil(t, lifecycle)
}
