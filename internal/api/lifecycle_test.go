// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api_test

import (
	"context"
	"errors"
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

// lifecyclePollClient returns a client polling lifecycleTestLocation. Authentication runs in test mode so that
// no real token is requested; with testMode false the client sleeps for real between polls.
func lifecyclePollClient(testMode bool) *api.Client {
	clientCfg := config.ProviderConfig{
		TestMode: testMode,
		Urls:     lifecycleTestUrls(),
	}
	authCfg := clientCfg
	authCfg.TestMode = true
	return api.NewApiClientBase(&clientCfg, api.NewAuthBase(&authCfg))
}

func lifecycleAcceptedResponse() *api.Response {
	return &api.Response{HttpResponse: &http.Response{Header: http.Header{
		constants.HEADER_LOCATION:    []string{lifecycleTestLocation},
		constants.HEADER_RETRY_AFTER: []string{"0"},
	}}}
}

// serveLifecycleOperation answers the polls of lifecycleTestLocation with responder through httpmock, which
// replaces the default transport until the test ends.
func serveLifecycleOperation(t *testing.T, responder httpmock.Responder) {
	t.Helper()
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	httpmock.RegisterResponder(http.MethodGet, lifecycleTestLocation, responder)
}

// errConnectionDropped is a poll that fails without an answer, which the client sees as a network error.
var errConnectionDropped = errors.New("connection reset by peer")

// hangUntilCanceled is a poll that never answers, until its request is canceled.
func hangUntilCanceled(r *http.Request) (*http.Response, error) {
	<-r.Context().Done()
	return nil, r.Context().Err()
}

func lifecycleJsonResponse(body string) *http.Response {
	resp := httpmock.NewStringResponse(http.StatusOK, body)
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

func TestUnitDoWaitForLifecycleOperationStatus_TransientPollFailureThenSuccess(t *testing.T) {
	callCount := 0
	serveLifecycleOperation(t, func(*http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return nil, errConnectionDropped
		}
		return lifecycleJsonResponse(`{"state":{"id":"Succeeded"}}`), nil
	})

	client := lifecyclePollClient(true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse())
	require.NoError(t, err, "a single failed poll must be retried on the next tick")
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_PersistentPollFailure_StopsAtConsecutiveFailureLimit(t *testing.T) {
	callCount := 0
	serveLifecycleOperation(t, func(*http.Request) (*http.Response, error) {
		callCount++
		return nil, errConnectionDropped
	})

	client := lifecyclePollClient(true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse())
	require.Error(t, err)
	require.ErrorContains(t, err, "failed 5 times in a row")
	var sentErr api.RequestSentError
	require.ErrorAs(t, err, &sentErr)
	require.Nil(t, lifecycle)
	require.Equal(t, constants.LIFECYCLE_POLL_MAX_CONSECUTIVE_FAILURES, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_PermanentPollFailure_FailsImmediately(t *testing.T) {
	callCount := 0
	serveLifecycleOperation(t, func(*http.Request) (*http.Response, error) {
		callCount++
		return httpmock.NewStringResponse(http.StatusBadRequest, "bad request"), nil
	})

	client := lifecyclePollClient(true)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(context.Background(), lifecycleAcceptedResponse())
	var statusErr customerrors.UnexpectedHttpStatusCodeError
	require.ErrorAs(t, err, &statusErr)
	require.Equal(t, http.StatusBadRequest, statusErr.StatusCode)
	require.Nil(t, lifecycle)
	require.Equal(t, 1, callCount)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextExpiresDuringHungPoll(t *testing.T) {
	serveLifecycleOperation(t, hangUntilCanceled)

	client := lifecyclePollClient(false)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, lifecycleAcceptedResponse())
	require.ErrorIs(t, err, context.DeadlineExceeded, "the deadline of the whole wait is final and must not be retried")
	require.Nil(t, lifecycle)
}

func TestUnitDoWaitForLifecycleOperationStatus_ContextTimeout_ReportsLastObservedOperation(t *testing.T) {
	serveLifecycleOperation(t, func(*http.Request) (*http.Response, error) {
		return lifecycleJsonResponse(`{"state":{"id":"Running"},"links":{"environment":{"path":"/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001"}}}`), nil
	})

	client := lifecyclePollClient(false)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	resp := lifecycleAcceptedResponse()
	resp.HttpResponse.Header.Set(constants.HEADER_RETRY_AFTER, "120")

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, resp)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, lifecycle)
	var waitErr *api.LifecycleWaitError
	require.ErrorAs(t, err, &waitErr)
	require.Equal(t, "Running", waitErr.LastObserved.State.Id)
	require.Equal(t, "/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001", waitErr.LastObserved.Links.Environment.Path)
}

// A Retry-After longer than the attempt deadline is a delay the service asked for, not a hung
// attempt, so it must be waited out under the deadline of the whole wait.
func TestUnitDoWaitForLifecycleOperationStatus_RetryAfterLongerThanAttemptTimeout_IsWaitedOut(t *testing.T) {
	callCount := 0
	serveLifecycleOperation(t, func(*http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			resp := httpmock.NewStringResponse(http.StatusTooManyRequests, "too many requests")
			resp.Header.Set(constants.HEADER_RETRY_AFTER, "1")
			return resp, nil
		}
		return lifecycleJsonResponse(`{"state":{"id":"Succeeded"}}`), nil
	})

	// Sleeps for real, so that the Retry-After delay outlasts the attempt deadline.
	client := lifecyclePollClient(false)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = api.WithRequestAttemptTimeout(ctx, 200*time.Millisecond)

	started := time.Now()
	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, lifecycleAcceptedResponse())
	require.NoError(t, err)
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount, "the 429 must be retried once, after the delay the service asked for")
	require.GreaterOrEqual(t, time.Since(started), time.Second, "the Retry-After delay must not be cut short")
}

func TestUnitDoWaitForLifecycleOperationStatus_HungAttempt_IsRetried(t *testing.T) {
	callCount := 0
	serveLifecycleOperation(t, func(r *http.Request) (*http.Response, error) {
		callCount++
		if callCount == 1 {
			return hangUntilCanceled(r)
		}
		return lifecycleJsonResponse(`{"state":{"id":"Succeeded"}}`), nil
	})

	client := lifecyclePollClient(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = api.WithRequestAttemptTimeout(ctx, 100*time.Millisecond)

	lifecycle, err := client.DoWaitForLifecycleOperationStatus(ctx, lifecycleAcceptedResponse())
	require.NoError(t, err, "a hung attempt must be cut by its own deadline and retried on the next tick")
	require.NotNil(t, lifecycle)
	require.Equal(t, "Succeeded", lifecycle.State.Id)
	require.Equal(t, 2, callCount)
}
