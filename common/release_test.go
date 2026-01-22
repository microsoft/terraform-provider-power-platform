// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package common

import "testing"

func TestReleaseDefaults(t *testing.T) {
	if ProviderVersion == "" {
		t.Fatal("expected ProviderVersion to be set")
	}
	if Commit == "" {
		t.Fatal("expected Commit to be set")
	}
	if Branch == "" {
		t.Fatal("expected Branch to be set")
	}
}
