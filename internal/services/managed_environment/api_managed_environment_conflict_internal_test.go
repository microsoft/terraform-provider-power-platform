// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

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

const managedEnvironmentConflictTestEnvironmentID = "00000000-0000-0000-0000-000000000001"

func TestUnitDisableManagedEnvironment_ConflictWhileStillManagedRetries(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := newManagedEnvironmentConflictTestClient()
	disableAttempts := 0
	environmentReads := 0

	registerDisableManagedEnvironmentResponder(t, func(_ *http.Request) (*http.Response, error) {
		disableAttempts++
		if disableAttempts == 1 {
			return httpmock.NewStringResponse(http.StatusConflict, ""), nil
		}
		return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
	})
	registerManagedEnvironmentStateResponder(t, func(_ *http.Request) (*http.Response, error) {
		environmentReads++
		return managedEnvironmentStateResponse(constants.PROTECTION_LEVEL_STANDARD), nil
	})

	err := client.DisableManagedEnvironment(context.Background(), managedEnvironmentConflictTestEnvironmentID)

	require.NoError(t, err)
	require.Equal(t, 2, disableAttempts, "a 409 while the environment is still managed must be retried")
	require.Equal(t, 1, environmentReads, "the provider must read the environment before deciding that a 409 is success")
}

func TestUnitDisableManagedEnvironment_ConflictWhenAlreadyDisabledIsIdempotent(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := newManagedEnvironmentConflictTestClient()
	disableAttempts := 0
	environmentReads := 0

	registerDisableManagedEnvironmentResponder(t, func(_ *http.Request) (*http.Response, error) {
		disableAttempts++
		return httpmock.NewStringResponse(http.StatusConflict, ""), nil
	})
	registerManagedEnvironmentStateResponder(t, func(_ *http.Request) (*http.Response, error) {
		environmentReads++
		return managedEnvironmentStateResponse("Basic"), nil
	})

	err := client.DisableManagedEnvironment(context.Background(), managedEnvironmentConflictTestEnvironmentID)

	require.NoError(t, err)
	require.Equal(t, 1, disableAttempts, "an already-disabled environment must not be retried")
	require.Equal(t, 1, environmentReads, "idempotent success must be based on the observed desired state")
}

func TestUnitDisableManagedEnvironment_PersistentConflictTerminatesWithoutFalseSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := newManagedEnvironmentConflictTestClient()
	disableAttempts := 0
	environmentReads := 0

	registerDisableManagedEnvironmentResponder(t, func(_ *http.Request) (*http.Response, error) {
		disableAttempts++
		return httpmock.NewStringResponse(http.StatusConflict, ""), nil
	})
	registerManagedEnvironmentStateResponder(t, func(_ *http.Request) (*http.Response, error) {
		environmentReads++
		return managedEnvironmentStateResponse(constants.PROTECTION_LEVEL_STANDARD), nil
	})

	err := client.DisableManagedEnvironment(context.Background(), managedEnvironmentConflictTestEnvironmentID)

	require.Error(t, err)
	require.ErrorContains(t, err, "maximum retries")
	require.Equal(t, constants.MAX_RETRY_COUNT+1, disableAttempts)
	require.Equal(t, constants.MAX_RETRY_COUNT+1, environmentReads)
}

func newManagedEnvironmentConflictTestClient() client {
	cfg := &config.ProviderConfig{
		TestMode: true,
		Urls: config.ProviderConfigUrls{
			BapiUrl:        "api.bap.microsoft.com",
			PowerAppsUrl:   "api.bap.microsoft.com",
			PowerAppsScope: "scope",
		},
	}
	return newManagedEnvironmentClient(api.NewApiClientBase(cfg, api.NewAuthBase(cfg)))
}

func registerDisableManagedEnvironmentResponder(t *testing.T, responder httpmock.Responder) {
	t.Helper()
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("no responder found for %s %s", req.Method, req.URL)
	})
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/"+managedEnvironmentConflictTestEnvironmentID+"/governanceConfiguration?api-version=2021-04-01",
		responder,
	)
}

func registerManagedEnvironmentStateResponder(t *testing.T, responder httpmock.Responder) {
	t.Helper()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/"+managedEnvironmentConflictTestEnvironmentID+"?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		responder,
	)
}

func managedEnvironmentStateResponse(protectionLevel string) *http.Response {
	return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
  "id": %q,
  "name": "env",
  "properties": {
    "governanceConfiguration": {"protectionLevel": %q}
  }
}`, managedEnvironmentConflictTestEnvironmentID, protectionLevel))
}
