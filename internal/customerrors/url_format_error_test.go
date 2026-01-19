// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import (
	"errors"
	"strings"
	"testing"
)

func TestUnitUrlFormatError(t *testing.T) {
	t.Parallel()

	wrapped := errors.New("parse failure")
	err := NewUrlFormatError("http://example", wrapped)

	ufe, ok := err.(*UrlFormatError)
	if !ok {
		t.Fatalf("expected UrlFormatError type")
	}

	if ufe.Url != "http://example" || !errors.Is(ufe, wrapped) {
		t.Fatalf("unexpected fields on UrlFormatError")
	}

	if got := ufe.Error(); !strings.Contains(got, "Request url must be an absolute url") || !strings.Contains(got, "parse failure") {
		t.Fatalf("unexpected error string: %q", got)
	}

	bare := NewUrlFormatError("http://example", nil)
	if got := bare.Error(); !strings.Contains(got, "http://example") {
		t.Fatalf("unexpected bare error string: %q", got)
	}
}
