// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customerrors

import "testing"

func TestUnitUnexpectedHttpStatusCodeError(t *testing.T) {
	t.Parallel()

	err := NewUnexpectedHttpStatusCodeError([]int{200, 202}, 500, "Internal Server Error", []byte("body"))
	httpErr, ok := err.(UnexpectedHttpStatusCodeError)
	if !ok {
		t.Fatal("expected UnexpectedHttpStatusCodeError type")
	}

	if httpErr.StatusCode != 500 || httpErr.StatusText != "Internal Server Error" {
		t.Fatalf("unexpected fields: %#v", httpErr)
	}

	if len(httpErr.ExpectedStatusCodes) != 2 || httpErr.ExpectedStatusCodes[0] != 200 || httpErr.ExpectedStatusCodes[1] != 202 {
		t.Fatalf("unexpected expected status codes: %#v", httpErr.ExpectedStatusCodes)
	}

	if string(httpErr.Body) != "body" {
		t.Fatalf("unexpected body: %q", string(httpErr.Body))
	}

	expected := "Unexpected HTTP status code. Expected: [200 202], received: [500] Internal Server Error | body"
	if got := httpErr.Error(); got != expected {
		t.Fatalf("unexpected error message: %q", got)
	}
}
