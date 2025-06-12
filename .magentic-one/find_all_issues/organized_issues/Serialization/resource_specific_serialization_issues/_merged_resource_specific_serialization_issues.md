# Resource-Specific Serialization Issues

This document consolidates issues related to serialization, struct field naming, and error handling in specific resource implementations.

## ISSUE 1

### Non-Idiomatic Client Struct Field Name

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go`

**Problem:** The field `EnvironmentWaveClient` in the `Resource` struct uses mixed casing (PascalCase) for a field that could be unexported. In Go, struct field names that do not need to be exported should start with a lowercase letter to indicate package-private scope and follow idiomatic Go conventions (unless the field is intended for serialization or use outside the current package).

**Impact:** Non-idiomatic naming makes the codebase less consistent with Go conventions and could lead to confusion or accidental exposure of private details of a struct. Overall severity: **medium**.

**Location:**

```go
// Line (approximately 18)
type Resource struct {
 helpers.TypeInfo
 EnvironmentWaveClient *environmentWaveClient
}
```

**Code Issue:**

```go
type Resource struct {
 helpers.TypeInfo
 EnvironmentWaveClient *environmentWaveClient
}
```

**Fix:** Make the field unexported by starting its name with a lowercase letter, unless it needs to be exported for reasons such as serialization, package usage, or framework requirements:

```go
type Resource struct {
 helpers.TypeInfo
 environmentWaveClient *environmentWaveClient
}
```

**Explanation:**

- This adheres to Go's idiomatic naming convention for struct fields.
- Only exported fields (uppercase) are accessible from outside the package.

## ISSUE 2

### Unreachable code or inconsistent state after removing all security roles in Read (environment user)

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

**Problem:** In the `Read` method (when handling an environment user), if all security roles are removed, a special case triggers assignment of `user.AadObjectId` and `user.DomainName` based on state, but the code then continues with marshalling and converting the (potentially incomplete) `user`.

**Impact:** This can result in setting state with partially populated user data, leading to data drift between Terraform and the actual environment, possibly breaking subsequent operations or plan consistency. Severity: **Medium**.

**Location:** From `Read` in "environment" branch:

```go
user, err := r.UserClient.GetEnvironmentUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
// if all the security roles are removed, the user will not be found
if user.AadObjectId == "" {
 user.AadObjectId = state.AadId.ValueString()
 user.DomainName = state.UserPrincipalName.ValueString()
}
if err != nil {
...
}
rolesBytes, err := json.Marshal(user.SecurityRoles)
...
updateUser = *user
```

**Code Issue:**

```go
if user.AadObjectId == "" {
 user.AadObjectId = state.AadId.ValueString()
 user.DomainName = state.UserPrincipalName.ValueString()
}
// (error handling for user == nil and err follows, possibly too late)
```

**Fix:** Move error handling and nil checks before attempting to operate on the `user` pointer. If the result indicates a removed user (all roles removed / not found), return early and clear state:

```go
user, err := r.UserClient.GetEnvironmentUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
if err != nil {
    if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
    return
}

if user == nil || user.AadObjectId == "" {
    // User not present, remove from state and exit
    resp.State.RemoveResource(ctx)
    return
}

// Only now safe to use user safely
rolesBytes, err := json.Marshal(user.SecurityRoles)
// ... etc
```

---

Apply this fix to the whole codebase

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
