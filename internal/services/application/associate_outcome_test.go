// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"errors"
	"net/http"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

// The association POST's classifier decides whether a failure proves the mutation was rejected or
// leaves the outcome unknown, and only errors the client marked as coming from the POST itself may
// classify at all. This pins the boundary per status and per phase.
func TestUnitIsAmbiguousAssociateFailure_Classifies_Every_Status(t *testing.T) {
	post := func(err error) error { return associationPostError{err: err} }
	statusErr := func(status int) error {
		return customerrors.NewUnexpectedHttpStatusCodeError([]int{http.StatusNoContent}, status, http.StatusText(status), nil)
	}

	definitive := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		// 425 Too Early is a refusal to process the request at all, so nothing was committed.
		http.StatusTooEarly,
	}
	for _, status := range definitive {
		if isAmbiguousAssociateFailure(post(statusErr(status))) {
			t.Errorf("status %d is a definitive rejection and must not be classified as ambiguous", status)
		}
	}

	ambiguous := []int{
		// An unexpected success or redirect is the opposite of a proven rejection.
		http.StatusOK,
		http.StatusAccepted,
		http.StatusMultipleChoices,
		http.StatusFound,
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
	for _, status := range ambiguous {
		if !isAmbiguousAssociateFailure(post(statusErr(status))) {
			t.Errorf("status %d hides the outcome and must be classified as ambiguous", status)
		}
	}

	if !isAmbiguousAssociateFailure(post(api.RequestSentError{Err: errors.New("connection reset by peer")})) {
		t.Error("a marked transport failure carries no status and must be classified as ambiguous")
	}
	if isAmbiguousAssociateFailure(api.RequestSentError{Err: errors.New("connection reset by peer")}) {
		t.Error("a transport failure without the POST marker happened before the association was attempted and must not be classified as ambiguous")
	}
	if isAmbiguousAssociateFailure(errors.New("could not resolve the environment host")) {
		t.Error("a plain pre-POST failure must not be classified as ambiguous")
	}
}
