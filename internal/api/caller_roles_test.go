// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The validation message Dataverse returned for a solution import into an organization created minutes earlier.
const callerWithoutRolesMessage = "TryGetMaxPrivilegeDepthForUserAcrossBusinessUnits: The user with id a068a392-04d1-f111-a059-000d3a374fad has not been assigned any roles. They need a role with the prvImportCustomization privilege."

func TestUnitIsCallerWithoutRolesMessage(t *testing.T) {
	assert.True(t, IsCallerWithoutRolesMessage(callerWithoutRolesMessage))
	assert.True(t, IsCallerWithoutRolesMessage("The user with id 1 HAS NOT BEEN ASSIGNED ANY ROLES."))
	assert.False(t, IsCallerWithoutRolesMessage("The solution requires a newer version of component X."))
}

func newCallerRolesTestClient() *Client {
	cfg := config.ProviderConfig{TestMode: true}
	return NewApiClientBase(&cfg, NewAuthBase(&cfg))
}

func TestUnitReplayWhileNewDataverseOrgAssignsCallerRoles_ReplaysUntilTheRolesArrive(t *testing.T) {
	const url = "https://caller-roles-arrive.crm.dynamics.com/api/data/v9.2/StageSolution"
	RecordNewDataverseOrg("https://caller-roles-arrive.crm.dynamics.com/")

	var logOutput bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &logOutput)
	runs := 0
	err := newCallerRolesTestClient().ReplayWhileNewDataverseOrgAssignsCallerRoles(ctx, url, func() (bool, error) {
		runs++
		if runs <= 3 {
			return true, errors.New(callerWithoutRolesMessage)
		}
		return false, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 4, runs)

	entries, decodeErr := tflogtest.MultilineJSONDecode(&logOutput)
	require.NoError(t, decodeErr)
	var windows []map[string]any
	for _, entry := range entries {
		if _, ok := entry["caller_without_roles_window"]; ok {
			windows = append(windows, entry)
		}
	}
	if assert.Len(t, windows, 1, "the end of the wait for the caller's roles must be logged once") {
		assert.Equal(t, "WARN", strings.ToUpper(fmt.Sprint(windows[0]["@level"])), "the nightly acceptance run logs at WARN, so the measurement must be a warning")
		assert.InDelta(t, 3, windows[0]["caller_without_roles_replays"], 0)
		assert.Equal(t, "was accepted", windows[0]["caller_without_roles_outcome"])
		assert.Equal(t, url, windows[0]["request_url_no_query"])
	}
}

func TestUnitReplayWhileNewDataverseOrgAssignsCallerRoles_StopsAtTheTimeBudget(t *testing.T) {
	RecordNewDataverseOrg("https://caller-roles-never.crm.dynamics.com/")

	runs := 0
	refusal := errors.New(callerWithoutRolesMessage)
	err := newCallerRolesTestClient().ReplayWhileNewDataverseOrgAssignsCallerRoles(context.Background(), "https://caller-roles-never.crm.dynamics.com/api/data/v9.2/StageSolution", func() (bool, error) {
		runs++
		return true, refusal
	})

	replays := int(constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_BUDGET / constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_INTERVAL)
	assert.Equal(t, replays+1, runs, "the replays must stop once their delays add up to the time budget")
	assert.ErrorIs(t, err, refusal)
	assert.ErrorContains(t, err, fmt.Sprintf("the caller still had no security roles after %d retries over %s", replays, constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_BUDGET))
}

func TestUnitReplayWhileNewDataverseOrgAssignsCallerRoles_DoesNotReplayAnOrganizationItDidNotCreate(t *testing.T) {
	runs := 0
	refusal := errors.New(callerWithoutRolesMessage)
	err := newCallerRolesTestClient().ReplayWhileNewDataverseOrgAssignsCallerRoles(context.Background(), "https://caller-roles-existing-org.crm.dynamics.com/api/data/v9.2/StageSolution", func() (bool, error) {
		runs++
		return true, refusal
	})

	assert.Equal(t, 1, runs, "in an organization the provider did not just create, a caller without roles is a configuration error")
	assert.Equal(t, refusal, err)
}

func TestUnitReplayWhileNewDataverseOrgAssignsCallerRoles_DoesNotReplayOtherErrors(t *testing.T) {
	RecordNewDataverseOrg("https://caller-roles-other-error.crm.dynamics.com/")

	runs := 0
	failure := errors.New("missing dependency: 'msdyn_Something'")
	err := newCallerRolesTestClient().ReplayWhileNewDataverseOrgAssignsCallerRoles(context.Background(), "https://caller-roles-other-error.crm.dynamics.com/api/data/v9.2/StageSolution", func() (bool, error) {
		runs++
		return false, failure
	})

	assert.Equal(t, 1, runs)
	assert.Equal(t, failure, err)
}
