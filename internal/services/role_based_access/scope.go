// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import "fmt"

// assignmentScope is where a role assignment applies. The RBAC API expresses scope twice for every
// call: once as a URL segment under /authorization, and once as a fully qualified string under the
// tenant. Both use the same segment, so it is built in one place here.
//
// The kind is always explicit, so an unset or empty id can never silently widen an assignment to
// the tenant.
type assignmentScope struct {
	kind string
	id   string
}

// Scope kinds, as written in the resource's `scope_type` attribute.
const (
	scopeTenant           = "tenant"
	scopeEnvironment      = "environment"
	scopeEnvironmentGroup = "environment_group"
)

// scopeKinds are the values the `scope_type` attribute accepts.
var scopeKinds = []string{scopeTenant, scopeEnvironment, scopeEnvironmentGroup}

func tenantAssignmentScope() assignmentScope {
	return assignmentScope{kind: scopeTenant}
}

func environmentAssignmentScope(environmentId string) assignmentScope {
	return assignmentScope{kind: scopeEnvironment, id: environmentId}
}

func environmentGroupAssignmentScope(environmentGroupId string) assignmentScope {
	return assignmentScope{kind: scopeEnvironmentGroup, id: environmentGroupId}
}

// segment is the path fragment shared by the request URL and the qualified scope string.
func (s assignmentScope) segment() string {
	switch s.kind {
	case scopeEnvironment:
		return fmt.Sprintf("/environments/%s", s.id)
	case scopeEnvironmentGroup:
		return fmt.Sprintf("/environmentGroups/%s", s.id)
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
	switch s.kind {
	case scopeEnvironment:
		return fmt.Sprintf("environment %s", s.id)
	case scopeEnvironmentGroup:
		return fmt.Sprintf("environment group %s", s.id)
	default:
		return "tenant"
	}
}
