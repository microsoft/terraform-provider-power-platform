// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestUnitProviderError_ErrorAndIs(t *testing.T) {
	t.Parallel()

	inner := errors.New("boom")
	err := ProviderError{ErrorCode: ErrorCode("CODE"), Err: inner}

	if got := err.Error(); got != "CODE: boom" {
		t.Fatalf("unexpected error string %q", got)
	}

	if !errors.Is(err, ProviderError{ErrorCode: "CODE"}) {
		t.Fatalf("expected errors.Is to match on error code")
	}

	if errors.Is(err, ProviderError{ErrorCode: "OTHER"}) {
		t.Fatalf("did not expect errors.Is to match different code")
	}
}

func TestUnitProviderError_UnwrapAndCode(t *testing.T) {
	t.Parallel()

	inner := errors.New("wrapped")
	wrapped := ProviderError{ErrorCode: ErrorCode("WRAPPED"), Err: fmt.Errorf("outer: %w", inner)}

	if got := Unwrap(wrapped); !errors.Is(got, inner) {
		t.Fatalf("unwrap did not return inner error")
	}

	if code := Code(wrapped); code != ErrorCode("WRAPPED") {
		t.Fatalf("expected code WRAPPED, got %s", code)
	}

	if code := Code(nil); code != "" {
		t.Fatalf("expected empty code for nil error, got %s", code)
	}
}

func TestUnitNewAndWrapProviderError(t *testing.T) {
	t.Parallel()

	base := NewProviderError(ErrorCode("BASE"), "message %d", 1)
	if got := base.Error(); !strings.Contains(got, "message 1") || !strings.Contains(got, "BASE") {
		t.Fatalf("expected formatted provider error, got %q", got)
	}

	wrapped := WrapIntoProviderError(errors.New("inner"), ErrorCode("WRAP"), "outer")
	if got := wrapped.Error(); got != "WRAP: outer: [inner]" {
		t.Fatalf("unexpected wrapped error string: %q", got)
	}

	nilWrapped := WrapIntoProviderError(nil, ErrorCode("WRAP"), "outer")
	if got := nilWrapped.Error(); got != "WRAP: outer" {
		t.Fatalf("unexpected nil wrapped error string: %q", got)
	}
}
