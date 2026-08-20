// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

// The create POST is not idempotent, so the classifier decides whether a failure is a definitive
// rejection (nothing committed, safe to report as-is) or leaves the outcome unknown (must become
// the explicit unknown-outcome error, with nothing adopted). This pins the boundary per status.
func TestUnitIsAmbiguousCreateFailure_Classifies_Every_Status(t *testing.T) {
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
		err := customerrors.NewUnexpectedHttpStatusCodeError([]int{http.StatusCreated}, status, http.StatusText(status), nil)
		if isAmbiguousCreateFailure(err) {
			t.Errorf("status %d is a definitive rejection and must not be classified as ambiguous", status)
		}
	}

	ambiguous := []int{
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
	for _, status := range ambiguous {
		err := customerrors.NewUnexpectedHttpStatusCodeError([]int{http.StatusCreated}, status, http.StatusText(status), nil)
		if !isAmbiguousCreateFailure(err) {
			t.Errorf("status %d hides the outcome and must be classified as ambiguous", status)
		}
	}

	if !isAmbiguousCreateFailure(errors.New("connection reset by peer")) {
		t.Error("a transport failure carries no status and must be classified as ambiguous")
	}
}

// The unknown-outcome error must carry the recovery instructions and keep the original failure
// reachable for callers that unwrap it.
func TestUnitUnknownOutcomeError_Wraps_And_Instructs(t *testing.T) {
	original := customerrors.NewUnexpectedHttpStatusCodeError([]int{http.StatusCreated}, http.StatusBadGateway, "Bad Gateway", nil)
	err := unknownOutcomeError(original)

	for _, phrase := range []string{
		"the create outcome is unknown",
		"no idempotency key or correlation identifier",
		"import it using the scope-shaped import id",
		"Terraform has not recorded an assignment in state",
	} {
		if !strings.Contains(err.Error(), phrase) {
			t.Errorf("the unknown-outcome error must contain %q, got: %s", phrase, err.Error())
		}
	}

	var httpErr customerrors.UnexpectedHttpStatusCodeError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadGateway {
		t.Error("the unknown-outcome error must wrap the original failure")
	}
}
