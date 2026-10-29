// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_isolation_policy //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/stretchr/testify/assert"
)

const operationNotFoundBody = `{"error":{"code":"OperationNotFound","message":"The operation '00000000-0000-0000-0000-000000000001' for object ID 'SERVICE' was not found."}}`

// pollingFixture is a client whose delays are recorded instead of slept, and the accepted response of a
// policy PUT whose operation location is served by a test server.
type pollingFixture struct {
	client   *Client
	accepted *api.Response
	sleeps   []time.Duration
	polls    int
}

// newPollingFixture serves the operation with answer, which receives the number of the poll.
func newPollingFixture(t *testing.T, answer func(poll int, w http.ResponseWriter)) *pollingFixture {
	t.Helper()
	f := &pollingFixture{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		f.polls++
		answer(f.polls, w)
	}))
	t.Cleanup(server.Close)

	cfg := config.ProviderConfig{TestMode: true}
	f.client = &Client{
		Api: api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg)),
		sleep: func(_ context.Context, duration time.Duration) error {
			f.sleeps = append(f.sleeps, duration)
			return nil
		},
	}
	f.accepted = &api.Response{HttpResponse: &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     http.Header{constants.HEADER_LOCATION: []string{server.URL + "/operations/00000000-0000-0000-0000-000000000001"}},
	}}
	return f
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, operationNotFoundBody)
}

// The operation service answered 404 OperationNotFound to the first poll of an operation the PUT had
// just accepted (nightly 2026-10-28), which used to fail the whole create.
func TestUnitTenantIsolationPolicyOperationPoll_Retries_OperationNotFound_Until_It_Appears(t *testing.T) {
	f := newPollingFixture(t, func(poll int, w http.ResponseWriter) {
		switch poll {
		case 1, 2:
			notFound(w)
		case 3:
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})

	_, err := f.client.doWaitForLifecycleOperationStatus(context.Background(), f.accepted)

	assert.NoError(t, err, "a 404 for an operation that appears shortly after must not fail the wait")
	assert.Equal(t, 4, f.polls)
	assert.Len(t, f.sleeps, 4, "the wait sleeps before the first poll and after every poll that is not final")
}

func TestUnitTenantIsolationPolicyOperationPoll_Gives_Up_On_OperationNotFound_After_Timeout(t *testing.T) {
	f := newPollingFixture(t, func(_ int, w http.ResponseWriter) {
		notFound(w)
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
	f := newPollingFixture(t, func(poll int, w http.ResponseWriter) {
		if poll <= 2 {
			w.Header().Set(constants.HEADER_RETRY_AFTER, "7")
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	_, err := f.client.doWaitForLifecycleOperationStatus(context.Background(), f.accepted)

	assert.NoError(t, err)
	assert.Equal(t, 3, f.polls)
	assert.Equal(t, []time.Duration{getRetryAfterDuration(nil), 7 * time.Second, 7 * time.Second}, f.sleeps, "every 409 must be followed by the delay the service asked for")
}
