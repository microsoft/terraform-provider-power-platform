# Security Issues Summary

This document contains all identified security issues in the Terraform Provider Power Platform codebase. These issues have been merged from individual analysis files to provide a comprehensive overview of security concerns that need to be addressed.

## ISSUE 1

### Title: Lack of Input Validation of the `location` Parameter

**File:** `/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go`

**Problem:**
The `location` parameter is inserted directly into the URL path using `fmt.Sprintf`. There is no checking or sanitization performed. If `location` is empty, malformed, or contains unexpected/slash/control characters, this can lead to malformed URLs and possibly security issues (e.g., path traversal), or functional bugs.

**Impact:**
**Medium**. Bugs or vulnerabilities can occur if unsanitized or user-generated input is passed.

**Location:**
Construction of the URL path:

```go
Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
```

**Fix:**
Validate the `location` parameter for empty string and improper characters before using it to build the URL. For example:

```go
if strings.TrimSpace(location) == "" {
    return currencies, fmt.Errorf("location parameter cannot be empty")
}

// Optionally, further sanitize or restrict allowed characters.
```

You might also want to URL-encode the `location` variable (if appropriate).

## ISSUE 2

### Title: Potential String Injection Using Raw IDs in HTTP Path

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go`

**Problem:**
In multiple functions (such as `GetDataverseUserBySystemUserId`, `UpdateDataverseUser`, `DeleteDataverseUser`, `RemoveDataverseSecurityRoles`, `AddDataverseSecurityRoles`), the `systemUserId` (or `roleId`) is interpolated directly into the HTTP path string without any URL path escaping. This can result in malformed URLs if the IDs contain unexpected or invalid URL path characters, or in edge cases, could be vulnerable to injection if the ID is not properly filtered (for example, a specially crafted ID could alter the intended HTTP request).

**Impact:**
Severity: Medium

This could cause failures if IDs contain special characters, and potentially be a vector for path injection or ambiguous logs. While the risk of controlled injection is somewhat low if IDs are always GUIDs, this is not explicitly enforced, and defensive programming is preferred.

**Location:**
Any place where code like this appears (e.g., in GetDataverseUserBySystemUserId):

**Code Issue:**

```go
Path: "/api/data/v9.2/systemusers(" + systemUserId + ")",
```

**Fix:**
Use `url.PathEscape` to safely encode path parameters when building URLs with variable user input.

```go
Path: "/api/data/v9.2/systemusers(" + url.PathEscape(systemUserId) + ")",
```

Repeat this fix for all paths that interpolate IDs in a similar manner. This ensures safe, predictable URL formation and protects against malformed requests or path confusion.

## ISSUE 3

### Title: Use of hard-coded string for magic UUIDs

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

**Problem:**
The schema field for `environment_routing_target_security_group_id` uses a hard-coded magic UUID value `00000000-0000-0000-0000-000000000000` in the description, but there is no reference to or enforcement for this value via constants. All code interacting with this field must treat this value specially, but the only place it is mentioned is in a schema description. This weakens data consistency and can lead to logic spread across the codebase.

**Impact:**
Hard to maintain, error-prone documentation and code; the expectation for a special UUID value is insulated in schema description only. If this logic is required elsewhere, it risks divergence and bugs. Severity: low.

**Location:**
Schema for `environment_routing_target_security_group_id`:

```go
"environment_routing_target_security_group_id": schema.StringAttribute{
    MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
```

**Fix:**
Declare a constant (if not already in use elsewhere) like:

```go
const ALLOW_ALL_USERS_UUID = "00000000-0000-0000-0000-000000000000"
```

Reference it in code and in descriptions by interpolating or documenting that constant. Additionally, logic that relies on this value should use this constant for data comparison and assignment, not inline string literals.

---

## Footer

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
