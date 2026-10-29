// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/stretchr/testify/require"
)

const (
	enterprisePolicyConflictTestEnvironmentID = "00000000-0000-0000-0000-000000000001"
	enterprisePolicyConflictTestSystemID      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
)

func TestUnitEnterprisePolicyOperation_ConflictUsesObservedDesiredState(t *testing.T) {
	tests := []struct {
		name             string
		action           string
		firstState       string
		wantPostAttempts int
	}{
		{
			name:             "link retries when the requested policy is absent",
			action:           "link",
			firstState:       enterprisePolicyEnvironmentWithoutPolicy(),
			wantPostAttempts: 2,
		},
		{
			name:             "link is idempotent when the requested policy is linked",
			action:           "link",
			firstState:       enterprisePolicyEnvironmentWithLinkedPolicy(),
			wantPostAttempts: 1,
		},
		{
			name:             "unlink retries while the requested policy remains linked",
			action:           "unlink",
			firstState:       enterprisePolicyEnvironmentWithLinkedPolicy(),
			wantPostAttempts: 2,
		},
		{
			name:             "unlink is idempotent when the requested policy is absent",
			action:           "unlink",
			firstState:       enterprisePolicyEnvironmentWithoutPolicy(),
			wantPostAttempts: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			client := newEnterprisePolicyConflictTestClient()
			postAttempts := 0
			environmentReads := 0

			httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
			})
			httpmock.RegisterResponder(
				http.MethodPost,
				client.buildEnterprisePolicyURL(enterprisePolicyConflictTestEnvironmentID, NETWORK_INJECTION_POLICY_TYPE, test.action),
				func(_ *http.Request) (*http.Response, error) {
					postAttempts++
					if postAttempts == 1 {
						return httpmock.NewStringResponse(http.StatusConflict, ""), nil
					}
					return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
				},
			)
			httpmock.RegisterResponder(
				http.MethodGet,
				"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+enterprisePolicyConflictTestEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
				func(_ *http.Request) (*http.Response, error) {
					environmentReads++
					return httpmock.NewStringResponse(http.StatusOK, test.firstState), nil
				},
			)

			err := client.executePolicyOperation(
				context.Background(),
				enterprisePolicyConflictTestEnvironmentID,
				NETWORK_INJECTION_POLICY_TYPE,
				enterprisePolicyConflictTestSystemID,
				test.action,
			)

			require.NoError(t, err)
			require.Equal(t, test.wantPostAttempts, postAttempts)
			require.Equal(t, 1, environmentReads, "a 409 may be success only after observing the requested desired state")
		})
	}
}

func TestUnitEnterprisePolicyOperation_PersistentConflictTerminatesWithoutFalseSuccess(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		remoteState string
	}{
		{
			name:        "link remains absent",
			action:      "link",
			remoteState: enterprisePolicyEnvironmentWithoutPolicy(),
		},
		{
			name:        "unlink remains linked",
			action:      "unlink",
			remoteState: enterprisePolicyEnvironmentWithLinkedPolicy(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			client := newEnterprisePolicyConflictTestClient()
			postAttempts := 0
			environmentReads := 0

			httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
			})
			httpmock.RegisterResponder(
				http.MethodPost,
				client.buildEnterprisePolicyURL(enterprisePolicyConflictTestEnvironmentID, NETWORK_INJECTION_POLICY_TYPE, test.action),
				func(_ *http.Request) (*http.Response, error) {
					postAttempts++
					return httpmock.NewStringResponse(http.StatusConflict, ""), nil
				},
			)
			httpmock.RegisterResponder(
				http.MethodGet,
				"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+enterprisePolicyConflictTestEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
				func(_ *http.Request) (*http.Response, error) {
					environmentReads++
					return httpmock.NewStringResponse(http.StatusOK, test.remoteState), nil
				},
			)

			err := client.executePolicyOperation(
				context.Background(),
				enterprisePolicyConflictTestEnvironmentID,
				NETWORK_INJECTION_POLICY_TYPE,
				enterprisePolicyConflictTestSystemID,
				test.action,
			)

			require.Error(t, err)
			require.ErrorContains(t, err, "maximum retries")
			require.Equal(t, constants.MAX_RETRY_COUNT+1, postAttempts)
			require.Equal(t, constants.MAX_RETRY_COUNT+1, environmentReads)
		})
	}
}

func newEnterprisePolicyConflictTestClient() Client {
	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl:        "api.bap.microsoft.com",
			PowerAppsUrl:   "api.bap.microsoft.com",
			PowerAppsScope: "scope",
		},
	}
	return newEnterprisePolicyClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
}

func enterprisePolicyEnvironmentWithoutPolicy() string {
	return fmt.Sprintf(`{
  "id": %q,
  "name": "env",
  "properties": {"enterprisePolicies": {}}
}`, enterprisePolicyConflictTestEnvironmentID)
}

func enterprisePolicyEnvironmentWithLinkedPolicy() string {
	return fmt.Sprintf(`{
  "id": %q,
  "name": "env",
  "properties": {
    "enterprisePolicies": {
      "vnets": {
        "systemId": %q,
        "linkStatus": "Linked"
      }
    }
  }
}`, enterprisePolicyConflictTestEnvironmentID, enterprisePolicyConflictTestSystemID)
}
