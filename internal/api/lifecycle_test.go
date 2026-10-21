// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&clientCfg))

	ctx := helpers.TestContext(context.Background(), t.Name())
	location := server.URL + "/lifecycle/status"
	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_OPERATION_LOCATION: []string{location},
		constants.HEADER_RETRY_AFTER:        []string{"0"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_SucceededNoSleep(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"state":{"id":"Succeeded"}}`))
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
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{server.URL + "/lifecycle/status"},
		constants.HEADER_RETRY_AFTER: []string{"120"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextTimeout(t *testing.T) {
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
	// Use TestMode=false for lifecycle polling, but keep auth in TestMode to avoid real token requests.
	authCfg.TestMode = true
	client := api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	resp := &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{server.URL + "/lifecycle/status"},
		constants.HEADER_RETRY_AFTER: []string{"120"},
	}}}

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, lifecycle)
}

// lifecyclePollClient returns a client polling server without sleeping between polls. Authentication
// runs in test mode so that no real token is requested.
func lifecyclePollClient(serverURL string, testMode bool) *api.Client {
	clientCfg := config.ProviderConfig{
		TestMode: testMode,
		Urls: config.ProviderConfigUrls{
			BapiUrl:            serverURL,
			PowerAppsUrl:       serverURL,
			PowerPlatformUrl:   serverURL,
			PowerAppsScope:     "scope",
			PowerPlatformScope: "scope",
		},
	}
	authCfg := clientCfg
	authCfg.TestMode = true
	return api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))
}

func lifecycleAcceptedResponse(serverURL string) *api.Response {
	return &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{serverURL + "/lifecycle/status"},
		constants.HEADER_RETRY_AFTER: []string{"0"},
	}}}
}

// dropConnection closes the connection without answering, which the client sees as a network error.
func dropConnection(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	conn, _, err := w.(http.Hijacker).Hijack()
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestUnitDoWaitForLifecycleOperationStatus_TransientPollFailureThenSuccess(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		if callCount == 1 {
			dropConnection(t, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"state":{"id":"Succeeded"}}`))
	}))
	defer server.Close()

	client := lifecyclePollClient(server.URL, true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse(server.URL))
	require.NoError(t, err, "a single failed poll must be retried on the next tick")
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_PersistentPollFailure_StopsAtConsecutiveFailureLimit(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		dropConnection(t, w)
	}))
	defer server.Close()

	client := lifecyclePollClient(server.URL, true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse(server.URL))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed 5 times in a row")
	var sentErr api.RequestSentError
	require.ErrorAs(t, err, &sentErr)
	require.Nil(t, lifecycle)
	require.Equal(t, constants.LIFECYCLE_POLL_MAX_CONSECUTIVE_FAILURES, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_PermanentPollFailure_FailsImmediately(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := lifecyclePollClient(server.URL, true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse(server.URL))
	var statusErr customerrors.UnexpectedHttpStatusCodeError
	require.ErrorAs(t, err, &statusErr)
	require.Equal(t, http.StatusBadRequest, statusErr.StatusCode)
	require.Nil(t, lifecycle)
	require.Equal(t, 1, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextExpiresDuringHungPoll(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := lifecyclePollClient(server.URL, false)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, lifecycleAcceptedResponse(server.URL))
	require.ErrorIs(t, err, context.DeadlineExceeded, "the deadline of the whole wait is final and must not be retried")
	require.Nil(t, lifecycle)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextTimeout_ReportsLastObservedOperation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"state":{"id":"Running"},"links":{"environment":{"path":"/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001"}}}`))
	}))
	defer server.Close()

	client := lifecyclePollClient(server.URL, false)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	resp := lifecycleAcceptedResponse(server.URL)
	resp.HttpResponse.Header.Set(constants.HEADER_RETRY_AFTER, "120")

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, lifecycle)
	var waitErr *api.LifecycleWaitError
	require.ErrorAs(t, err, &waitErr)
	require.Equal(t, "Running", waitErr.LastObserved.State.Id)
	require.Equal(t, "/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001", waitErr.LastObserved.Links.Environment.Path)
}
