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
		t.Fatal("expected errors.Is to match on error code")
	}

	if errors.Is(err, ProviderError{ErrorCode: "OTHER"}) {
		t.Fatal("did not expect errors.Is to match different code")
	}

	noInner := ProviderError{ErrorCode: ErrorCode("NO_INNER")}
	if got := noInner.Error(); got != "NO_INNER" {
		t.Fatalf("unexpected error string for nil inner: %q", got)
	}
}

func TestUnitProviderError_UnwrapAndCode(t *testing.T) {
	t.Parallel()

	inner := errors.New("wrapped")
	wrapped := ProviderError{ErrorCode: ErrorCode("WRAPPED"), Err: fmt.Errorf("outer: %w", inner)}

	if got := Unwrap(wrapped); !errors.Is(got, inner) {
		t.Fatal("unwrap did not return inner error")
	}

	if code := Code(wrapped); code != ErrorCode("WRAPPED") {
		t.Fatalf("expected code WRAPPED, got %s", code)
	}

	if code := Code(nil); code != "" {
		t.Fatalf("expected empty code for nil error, got %s", code)
	}

	if got := Unwrap(ProviderError{ErrorCode: ErrorCode("NO_INNER")}); got != nil {
		t.Fatalf("expected nil unwrap for ProviderError without inner error, got %v", got)
	}

	plainWrapped := fmt.Errorf("plain: %w", inner)
	if got := Unwrap(plainWrapped); !errors.Is(got, inner) {
		t.Fatal("expected unwrap to return inner error for non-provider error")
	}

	if code := Code(errors.New("plain")); code != "" {
		t.Fatalf("expected empty code for non-provider error, got %s", code)
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
