// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitBuildEnvironmentHostUri(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		environmentId string
		want          string
		wantErr       bool
	}{
		{name: "valid guid", environmentId: "00000000-0000-0000-0000-000000000123", want: "000000000000000000000000000001.23.environment.api.powerplatform.com"},
		{name: "empty", environmentId: "", wantErr: true},
		{name: "one character", environmentId: "a", wantErr: true},
		{name: "only dashes", environmentId: "---", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := helpers.BuildEnvironmentHostUri(tt.environmentId, "api.powerplatform.com")
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildEnvironmentHostUri() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("BuildEnvironmentHostUri() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnitBuildTenantHostUri(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tenantId string
		want     string
		wantErr  bool
	}{
		{name: "valid guid", tenantId: "00000000-0000-0000-0000-000000000123", want: "000000000000000000000000000001.23.tenant.api.powerplatform.com"},
		{name: "empty", tenantId: "", wantErr: true},
		{name: "one character", tenantId: "a", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := helpers.BuildTenantHostUri(tt.tenantId, "api.powerplatform.com")
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildTenantHostUri() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("BuildTenantHostUri() = %q, want %q", got, tt.want)
			}
		})
	}
}
