# Title

Missing documentation for structs and their fields.

## Path

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/models.go

## Problem

The file lacks documentation comments for the `Resource` and `sourceModel` structs, as well as their fields. Documentation provides clarity on the purpose of each struct and its fields, especially when working within collaborative teams or open-source projects.

## Impact

Without structured documentation:
1. Developers may struggle to understand the purpose of each struct and its fields, reducing readability/maintainability of the code.
2. It becomes harder for newcomers to the codebase to work effectively.
3. The integration with tools like Terraform or Go documentation generators may miss expected metadata.

Severity: **Low**

## Location

The issue is found in the declarations of the following:
1. Struct `Resource`.
2. Struct `sourceModel`.

## Code Issue

Here is the relevant code snippet:

```go
type Resource struct {
	helpers.TypeInfo
	EnterprisePolicyClient Client
}

type sourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	SystemId      types.String   `tfsdk:"system_id"`
	PolicyType    types.String   `tfsdk:"policy_type"`
}
```

## Fix

Add comments documenting these structs and fields, explaining what they represent and how they are meant to be used.

```go
// Resource represents a terraform-compatible resource used to 
// manage Enterprise Policies in Power Platform.
type Resource struct {
	// TypeInfo provides generic type metadata required by the terraform framework.
	helpers.TypeInfo
	// EnterprisePolicyClient is the client used to communicate with the Power Platform API.
	EnterprisePolicyClient Client
}

// sourceModel represents the values that can be configured or extracted 
// from an Enterprise Policy resource during terraform operations.
type sourceModel struct {
	// Timeouts specifies the timeout configurations for terraform operations.
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	// Id is the unique identifier of the Enterprise Policy resource.
	Id types.String `tfsdk:"id"`
	// EnvironmentId is the identifier of the environment this policy belongs to.
	EnvironmentId types.String `tfsdk:"environment_id"`
	// SystemId is the identifier of the system for which the policy applies.
	SystemId types.String `tfsdk:"system_id"`
	// PolicyType specifies the type of the policy being managed.
	PolicyType types.String `tfsdk:"policy_type"`
}
```