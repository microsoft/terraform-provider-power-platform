// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import "testing"

func TestUnitUnexpectedHttpStatusCodeError(t *testing.T) {
	t.Parallel()

	err := NewUnexpectedHttpStatusCodeError([]int{200, 202}, 500, "Internal Server Error", []byte("body"))
	httpErr, ok := err.(UnexpectedHttpStatusCodeError)
	if !ok {
		t.Fatalf("expected UnexpectedHttpStatusCodeError type")
	}

	if httpErr.StatusCode != 500 || httpErr.StatusText != "Internal Server Error" {
		t.Fatalf("unexpected fields: %#v", httpErr)
	}

	expected := "Unexpected HTTP status code. Expected: [200 202], received: [500] Internal Server Error | body"
	if got := httpErr.Error(); got != expected {
		t.Fatalf("unexpected error message: %q", got)
	}
}
