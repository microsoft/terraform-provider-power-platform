// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import "fmt"

// roleHolder is the Dataverse principal a security role is assigned to. Dataverse keeps users and
// teams in separate tables, each with its own role association, so the holder decides both the
// entity set to address and the association to write through.
type roleHolder struct {
	systemUserId string
	teamId       string
}

func systemUserRoleHolder(systemUserId string) roleHolder {
	return roleHolder{systemUserId: systemUserId}
}

func teamRoleHolder(teamId string) roleHolder {
	return roleHolder{teamId: teamId}
}

func (h roleHolder) isTeam() bool {
	return h.teamId != ""
}

func (h roleHolder) id() string {
	if h.isTeam() {
		return h.teamId
	}
	return h.systemUserId
}

// entitySet is the Dataverse collection holding this principal.
func (h roleHolder) entitySet() string {
	if h.isTeam() {
		return "teams"
	}
	return "systemusers"
}

// association is the navigation property linking this principal to its security roles.
func (h roleHolder) association() string {
	if h.isTeam() {
		return "teamroles_association"
	}
	return "systemuserroles_association"
}

// selectFields are the columns needed to resolve the principal's business unit and confirm it exists.
func (h roleHolder) selectFields() string {
	if h.isTeam() {
		return "teamid,name,_businessunitid_value"
	}
	return "applicationid,systemuserid,fullname,isdisabled,deletedstate,_businessunitid_value"
}

// path addresses the principal itself.
func (h roleHolder) path(apiVersion string) string {
	return fmt.Sprintf("/api/data/%s/%s(%s)", apiVersion, h.entitySet(), h.id())
}

// associationPath addresses the principal's security role collection.
func (h roleHolder) associationPath(apiVersion string) string {
	return fmt.Sprintf("%s/%s/$ref", h.path(apiVersion), h.association())
}

func (h roleHolder) String() string {
	if h.isTeam() {
		return fmt.Sprintf("team %s", h.teamId)
	}
	return fmt.Sprintf("system user %s", h.systemUserId)
}
