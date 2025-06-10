# Input Validation and Type Safety Issues

This document contains merged issues related to input validation and type safety in the Power Platform Terraform provider.

## ISSUE 1

**Title:** Lack of input validation for critical fields

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

**Problem:**
Some string attributes that represent UUID values (such as `environment_routing_target_environment_group_id` and `environment_routing_target_security_group_id`) do not have format validators to ensure valid UUIDs are supplied.

**Impact:**
Invalid data can be accepted at plan/apply time, and may lead to run-time errors or unexpected behavior interacting with the API. Data consistency and user experience are negatively affected. Severity: medium.

**Location:**
Schema definition for fields with `CustomType: customtypes.UUIDType{}` (e.g., under `.governance` nested field).

**Code Issue:**

```go
"environment_routing_target_environment_group_id": schema.StringAttribute{
    MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
"environment_routing_target_security_group_id": schema.StringAttribute{
    MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
```

**Fix:**
Add UUID format validators. For example:

```go
import "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

...

"environment_routing_target_environment_group_id": schema.StringAttribute{
    MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"),
            "Must be a valid UUID",
        ),
    },
},
```

And similarly for other UUID fields.

## ISSUE 2

**Title:** Lack of UUID format validation on NewUUIDValue

**File:** `/workspaces/terraform-provider-power-platform/internal/customtypes/uuid.go`

**Problem:**
The function `NewUUIDValue` accepts any string without validating whether the provided value is a valid UUID format. This can potentially allow invalid UUIDs to enter the system, which may lead to bugs or data inconsistencies later.

**Impact:**
This is a medium-severity type safety and data consistency issue. Storing invalid UUIDs can cause downstream errors, integration issues, or data corruption.

**Location:**

```go
func NewUUIDValue(value string) UUID {
 return UUID{
  StringValue: basetypes.NewStringValue(value),
 }
}
```

**Code Issue:**

```go
func NewUUIDValue(value string) UUID {
 return UUID{
  StringValue: basetypes.NewStringValue(value),
 }
}
```

**Fix:**
Validate the UUID format using a regular expression or a UUID parsing library before allowing the value. You may return an error or handle diagnostics if the value is invalid.

```go
import (
 "regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func NewUUIDValue(value string) UUID {
 if !uuidRegex.MatchString(value) {
  // handle as you see fit, you could return a Null/Unknown UUID
  return UUID{
   StringValue: basetypes.NewStringUnknown(),
  }
 }
 return UUID{
  StringValue: basetypes.NewStringValue(value),
 }
}
```

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
