// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy //nolint:revive // the underscored package name predates this file and matches every service in the repo

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
	"github.com/stretchr/testify/assert"
)

const operationNotFoundBody = `{"error":{"code":"OperationNotFound","message":"The operation '00000000-0000-0000-0000-000000000001' for object ID 'SERVICE' was not found."}}`

// operationUrl is an https URL on the BAPI host, because a request without explicit scopes is only
// authorized for https URLs of known hosts.
const operationUrl = "https://" + constants.PUBLIC_BAPI_DOMAIN + "/providers/PowerPlatform.Governance/v1/operations/00000000-0000-0000-0000-000000000001"

// pollingFixture is a client whose delays are recorded instead of slept, and the accepted response of a
// policy PUT whose operation httpmock serves.
type pollingFixture struct {
	client   *Client
	accepted *api.Response
	sleeps   []time.Duration
	polls    int
}

// newPollingFixture serves the operation with answer, which receives the number of the poll.
func newPollingFixture(t *testing.T, answer func(poll int) *http.Response) *pollingFixture {
	t.Helper()
	f := &pollingFixture{}
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	httpmock.RegisterResponder(http.MethodGet, operationUrl, func(*http.Request) (*http.Response, error) {
		f.polls++
		return answer(f.polls), nil
	})

	cfg := config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl:        constants.PUBLIC_BAPI_DOMAIN,
			PowerAppsScope: constants.PUBLIC_POWERAPPS_SCOPE,
		},
	}
	f.client = &Client{
		Api: api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg)),
		sleep: func(_ context.Context, duration time.Duration) error {
			f.sleeps = append(f.sleeps, duration)
			return nil
		},
	}
	f.accepted = &api.Response{HttpResponse: &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     http.Header{constants.HEADER_LOCATION: []string{operationUrl}},
	}}
	return f
}

func notFound() *http.Response {
	resp := httpmock.NewStringResponse(http.StatusNotFound, operationNotFoundBody)
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// The operation service answered 404 OperationNotFound to the first poll of an operation the PUT had
// just accepted (nightly 2026-10-28), which used to fail the whole create.
func TestUnitTenantIsolationPolicyOperationPoll_Retries_OperationNotFound_Until_It_Appears(t *testing.T) {
	f := newPollingFixture(t, func(poll int) *http.Response {
		switch poll {
		case 1, 2:
			return notFound()
		case 3:
			return httpmock.NewStringResponse(http.StatusAccepted, "")
		default:
			return httpmock.NewStringResponse(http.StatusOK, "")
		}
	})

	_, err := f.client.doWaitForLifecycleOperationStatus(context.Background(), f.accepted)

	assert.NoError(t, err, "a 404 for an operation that appears shortly after must not fail the wait")
	assert.Equal(t, 4, f.polls)
	assert.Len(t, f.sleeps, 4, "the wait sleeps before the first poll and after every poll that is not final")
}

func TestUnitTenantIsolationPolicyOperationPoll_Gives_Up_On_OperationNotFound_After_Timeout(t *testing.T) {
	f := newPollingFixture(t, func(int) *http.Response {
		return notFound()
	})

	_, err := f.client.doWaitForLifecycleOperationStatus(context.Background(), f.accepted)

	defaultDelay := getRetryAfterDuration(nil)
	assert.Equal(t, int(constants.TENANT_ISOLATION_OPERATION_NOT_FOUND_TIMEOUT/defaultDelay)+1, f.polls, "an operation that never appears must stop being polled once the 404 time budget is spent")
	assert.ErrorContains(t, err, "OperationNotFound")
	var statusErr customerrors.UnexpectedHttpStatusCodeError
	if assert.ErrorAs(t, err, &statusErr) {
		assert.Equal(t, http.StatusNotFound, statusErr.StatusCode)
	}
}

// A 409 used to be polled again at once, in a loop with no delay at all.
func TestUnitTenantIsolationPolicyOperationPoll_Waits_Between_Conflict_Polls(t *testing.T) {
	f := newPollingFixture(t, func(poll int) *http.Response {
		if poll <= 2 {
			resp := httpmock.NewStringResponse(http.StatusConflict, "")
			resp.Header.Set(constants.HEADER_RETRY_AFTER, "7")
			return resp
		}
		return httpmock.NewStringResponse(http.StatusOK, "")
	})

	_, err := f.client.doWaitForLifecycleOperationStatus(context.Background(), f.accepted)

	assert.NoError(t, err)
	assert.Equal(t, 3, f.polls)
	assert.Equal(t, []time.Duration{getRetryAfterDuration(nil), 7 * time.Second, 7 * time.Second}, f.sleeps, "every 409 must be followed by the delay the service asked for")
}
