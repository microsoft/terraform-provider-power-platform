// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type LifecycleDto struct {
	Id                 string                  `json:"id"`
	Links              LifecycleLinksDto       `json:"links"`
	State              LifecycleStateDto       `json:"state"`
	Type               LifecycleStateDto       `json:"type"`
	CreatedDateTime    string                  `json:"createdDateTime"`
	LastActionDateTime string                  `json:"lastActionDateTime"`
	RequestedBy        LifecycleRequestedByDto `json:"requestedBy"`
	Stages             []LifecycleStageDto     `json:"stages"`
}

type LifecycleStageDto struct {
	Id                  string            `json:"id"`
	Name                string            `json:"name"`
	State               LifecycleStateDto `json:"state"`
	FirstActionDateTime string            `json:"firstActionDateTime"`
	LastActionDateTime  string            `json:"lastActionDateTime"`
}

type LifecycleLinksDto struct {
	Self        LifecycleLinkDto `json:"self"`
	Environment LifecycleLinkDto `json:"environment"`
}

type LifecycleLinkDto struct {
	Path string `json:"path"`
}

type LifecycleStateDto struct {
	Id string `json:"id"`
}

type LifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

// LifecycleWaitError is returned when waiting for a lifecycle operation stopped before the operation
// reached a terminal state. LastObserved is the latest state the operation reported, so a caller
// can still identify the resource that an unfinished operation is creating. It adds no text of its
// own, so the wrapped error reads exactly as it did before.
type LifecycleWaitError struct {
	LastObserved *LifecycleDto
	Err          error
}

func (e *LifecycleWaitError) Error() string {
	return e.Err.Error()
}

func (e *LifecycleWaitError) Unwrap() error {
	return e.Err
}

// DoWaitForLifecycleOperationStatus polls the lifecycle operation referenced by the Location or
// Operation-Location header of response until it succeeds or fails. A poll that fails transiently is
// retried on the next tick, so only the deadline of ctx, a permanent error, or a run of
// constants.LIFECYCLE_POLL_MAX_CONSECUTIVE_FAILURES transient failures ends the wait early.
func (client *Client) DoWaitForLifecycleOperationStatus(ctx context.Context, response *Response) (*LifecycleDto, error) {
	locationHeader := response.GetHeader(constants.HEADER_LOCATION)
	if locationHeader == "" {
		locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION)
	}
	tflog.Debug(ctx, "Location Header: "+locationHeader)

	if locationHeader == "" {
		return nil, nil
	}

	_, err := url.Parse(locationHeader)
	if err != nil {
		tflog.Error(ctx, "Error parsing location header: "+err.Error())
	}

	waitFor := retryAfter(ctx, response.HttpResponse)

	var lastObserved *LifecycleDto
	waitError := func(err error) error {
		if lastObserved == nil {
			return err
		}
		return &LifecycleWaitError{LastObserved: lastObserved, Err: err}
	}

	consecutiveFailures := 0
	for {
		lifecycleResponse, response, err := client.pollLifecycleOperation(ctx, locationHeader)
		if err != nil {
			if ctx.Err() != nil || !isTransientLifecyclePollError(err) {
				return nil, waitError(err)
			}
			consecutiveFailures++
			if consecutiveFailures >= constants.LIFECYCLE_POLL_MAX_CONSECUTIVE_FAILURES {
				return nil, waitError(fmt.Errorf("lifecycle operation poll failed %d times in a row: %w", consecutiveFailures, err))
			}
			tflog.Warn(ctx, fmt.Sprintf("Lifecycle operation poll failed transiently (%d in a row), retrying: %s", consecutiveFailures, err.Error()))
			if err := client.SleepWithContext(ctx, waitFor); err != nil {
				return nil, waitError(err)
			}
			continue
		}
		consecutiveFailures = 0

		tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")

		if response.HttpResponse.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Lifecycle operation returned 404 - resource already gone, treating as success")
			// Return a non-nil LifecycleDto sentinel so callers can safely dereference the result.
			lifecycleResponse.State.Id = "Succeeded"
			return lifecycleResponse, nil
		}

		if response.HttpResponse.StatusCode == http.StatusConflict {
			err = client.SleepWithContext(ctx, waitFor)
			if err != nil {
				return nil, waitError(err)
			}
			continue
		}
		lastObserved = lifecycleResponse
		tflog.Debug(ctx, "Lifecycle Operation State: '"+lifecycleResponse.State.Id+"'")
		if lifecycleResponse.State.Id == "Succeeded" || lifecycleResponse.State.Id == "Failed" {
			return lifecycleResponse, nil
		}

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return nil, waitError(err)
		}
	}
}

// pollLifecycleOperation reads the lifecycle operation once. The read gets its own deadline because
// the HTTP client has none, so a single hung connection would otherwise consume the whole operation
// budget instead of being retried on the next tick.
func (client *Client) pollLifecycleOperation(ctx context.Context, location string) (*LifecycleDto, *Response, error) {
	pollCtx, cancel := context.WithTimeout(ctx, constants.LIFECYCLE_POLL_REQUEST_TIMEOUT)
	defer cancel()

	lifecycleResponse := LifecycleDto{}
	response, err := client.Execute(pollCtx, nil, "GET", location, nil, nil, []int{http.StatusOK, http.StatusConflict, http.StatusNotFound}, &lifecycleResponse)
	if err != nil && pollCtx.Err() != nil && ctx.Err() == nil {
		return nil, response, fmt.Errorf("lifecycle operation poll did not complete within %s: %w", constants.LIFECYCLE_POLL_REQUEST_TIMEOUT, err)
	}
	return &lifecycleResponse, response, err
}

// isTransientLifecyclePollError reports whether a failed lifecycle operation read is worth repeating.
// Reading the operation is idempotent, so a lost connection, a missing or unreadable response, a
// read that ran out of its own deadline, or a server-side status already treated as retryable
// says nothing about the operation itself. A refused credential or any other status is permanent.
func isTransientLifecyclePollError(err error) bool {
	var statusErr customerrors.UnexpectedHttpStatusCodeError
	if errors.As(err, &statusErr) {
		return statusErr.StatusCode != http.StatusUnauthorized && helpers.ArrayContains(retryableStatusCodes, statusErr.StatusCode)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		// The caller has already ruled out the deadline of the whole wait, so this is the read's own.
		return true
	}
	var sentErr RequestSentError
	return errors.As(err, &sentErr)
}
