// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package common

import "testing"

func TestReleaseDefaults(t *testing.T) {
	if ProviderVersion != "0.0.0-dev" {
		t.Fatalf("expected ProviderVersion default %q, got %q", "0.0.0-dev", ProviderVersion)
	}
	if Commit != "dev" {
		t.Fatalf("expected Commit default %q, got %q", "dev", Commit)
	}
	if Branch != "dev" {
		t.Fatalf("expected Branch default %q, got %q", "dev", Branch)
	}
}
