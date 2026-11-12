// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"testing"
	"time"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestUnitNewDataverseOrg_MatchesTheOrganizationHostOnly(t *testing.T) {
	RecordNewDataverseOrg("https://New-Org-Host-Test.crm4.dynamics.com/")

	assert.True(t, IsNewDataverseOrg("https://new-org-host-test.crm4.dynamics.com/api/data/v9.2/contacts?$top=1"), "requests to the organization's host match regardless of path, query and case")
	assert.False(t, IsNewDataverseOrg("https://other-org-host-test.crm4.dynamics.com/api/data/v9.2/contacts"))
	assert.False(t, IsNewDataverseOrg("::not a url"))
}

func TestUnitNewDataverseOrg_ExpiresAfterTheNewOrgPeriod(t *testing.T) {
	const host = "expired-org-host-test.crm.dynamics.com"
	RecordNewDataverseOrg("https://" + host + "/")
	assert.True(t, IsNewDataverseOrg("https://"+host+"/"))

	newDataverseOrgs.Lock()
	newDataverseOrgs.createdAt[host] = time.Now().Add(-constants.NEW_DATAVERSE_ORG_PERIOD - time.Minute)
	newDataverseOrgs.Unlock()

	assert.False(t, IsNewDataverseOrg("https://"+host+"/"), "an organization created longer ago than the period gets the ordinary 401 cap")

	RecordNewDataverseOrg("https://another-org-host-test.crm.dynamics.com/")
	newDataverseOrgs.Lock()
	_, stillKnown := newDataverseOrgs.createdAt[host]
	newDataverseOrgs.Unlock()
	assert.False(t, stillKnown, "recording an organization prunes the expired ones")
}
