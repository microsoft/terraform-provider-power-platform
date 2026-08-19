// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import "fmt"

// assignmentScope is where a role assignment applies. The RBAC API expresses scope twice for every
// call: once as a URL segment under /authorization, and once as a fully qualified string under the
// tenant. Both use the same segment, so it is built in one place here.
//
// An empty scope is the tenant itself.
type assignmentScope struct {
	environmentId      string
	environmentGroupId string
}

func tenantAssignmentScope() assignmentScope {
	return assignmentScope{}
}

func environmentAssignmentScope(environmentId string) assignmentScope {
	return assignmentScope{environmentId: environmentId}
}

func environmentGroupAssignmentScope(environmentGroupId string) assignmentScope {
	return assignmentScope{environmentGroupId: environmentGroupId}
}

// segment is the path fragment shared by the request URL and the qualified scope string.
func (s assignmentScope) segment() string {
	switch {
	case s.environmentId != "":
		return fmt.Sprintf("/environments/%s", s.environmentId)
	case s.environmentGroupId != "":
		return fmt.Sprintf("/environmentGroups/%s", s.environmentGroupId)
	default:
		return ""
	}
}

// collectionPath is the roleAssignments collection for this scope.
func (s assignmentScope) collectionPath() string {
	return fmt.Sprintf("/authorization%s/roleAssignments", s.segment())
}

// assignmentPath addresses a single role assignment within this scope.
func (s assignmentScope) assignmentPath(roleAssignmentId string) string {
	return fmt.Sprintf("%s/%s", s.collectionPath(), roleAssignmentId)
}

// qualify renders the scope string the API stores on the assignment.
func (s assignmentScope) qualify(tenantScope string) string {
	return tenantScope + s.segment()
}

// String describes the scope for logs and error messages.
func (s assignmentScope) String() string {
	switch {
	case s.environmentId != "":
		return fmt.Sprintf("environment %s", s.environmentId)
	case s.environmentGroupId != "":
		return fmt.Sprintf("environment group %s", s.environmentGroupId)
	default:
		return "tenant"
	}
}
