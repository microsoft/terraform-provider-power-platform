# JSON Tagging Issues

This document consolidates all issues related to missing, incorrect, or inconsistent JSON struct tags across the codebase.

## ISSUE 1

### Missing JSON Tags on Public Field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go`

**Problem:** The only field without a JSON tag in this file appears to be `InstanceURL` in the `linkedEnvironmentIdMetadataDto` struct. All other struct fields are explicitly tagged for marshaling/unmarshaling. Omitting a tag risks inconsistent casing/behavior for JSON consumers.

**Impact:** May cause issues with (de)serialization or integration, especially if relying on standard lowerCamelCase-to-snake_case auto-conversion, and is inconsistent with the rest of the DTO file. Severity: low.

**Location:** Line 127

**Code Issue:**

```go
type linkedEnvironmentIdMetadataDto struct {
 InstanceURL string
}
```

**Fix:** Add a JSON tag for the field, e.g.:

```go
type linkedEnvironmentIdMetadataDto struct {
 InstanceURL string `json:"instanceUrl"`
}
```

## ISSUE 2

### Incorrect JSON Tag for Struct Field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/dto.go`

**Problem:** The struct `linkedEnvironmentIdMetadataDto` defines a field `InstanceURL` without a JSON tag. All other struct fields in this file use explicit JSON tags to map Go struct fields to the correct JSON keys, likely for un/marshaling purposes from API responses or requests.

**Impact:** This oversight could result in unexpected JSON key casing (e.g., "InstanceURL" instead of "instanceUrl") when serializing or deserializing JSON, which can cause bugs in API communication and data inconsistencies. This is a **medium** severity issue due to its potential to cause subtle bugs in data exchange.

**Location:**

```go
type linkedEnvironmentIdMetadataDto struct {
 InstanceURL string
}
```

**Code Issue:**

```go
type linkedEnvironmentIdMetadataDto struct {
 InstanceURL string
}
```

**Fix:** Specify the correct JSON tag for the field, matching the expected API field name casing:

```go
type linkedEnvironmentIdMetadataDto struct {
 InstanceURL string `json:"instanceUrl"`
}
```

## ISSUE 3

### Missing JSON struct tag for Unblockable field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go`

**Problem:** The `Unblockable` field in `connectorPropertiesDto` lacks a JSON struct tag. This can lead to inconsistent marshaling and unmarshaling behaviors when working with JSON, potentially causing bugs if this struct is used with JSON APIs.

**Impact:** Medium severity: Omitting the JSON tag makes this field invisible when serializing or deserializing, which can lead to subtle bugs in API integration or data storage.

**Location:**

```go
type connectorPropertiesDto struct {
 DisplayName string `json:"displayName"`
 Description string `json:"description"`
 Tier        string `json:"tier"`
 Publisher   string `json:"publisher"`
 Unblockable bool
}
```

**Code Issue:**

```go
 Unblockable bool
```

**Fix:** Add an appropriate JSON struct tag to the field:

```go
 Unblockable bool `json:"unblockable"`
```

## ISSUE 4

### Boolean Field with `omitempty` Tag

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/dto.go`

**Problem:** Several boolean fields are tagged with `omitempty` in their JSON struct tags (e.g., `BingChatEnabled bool`json:"bingChatEnabled,omitempty"``). In Go, a boolean's zero value is `false`, and when using `omitempty`, a `false` value omits the field from encoded JSON. This can cause unintentional absence of the field, which could be ambiguous for API consumers.

**Impact:** Leads to ambiguity between `false` (explicit) and the field not being set at all, especially when the DTO evolves or is consumed by other systems. This is generally a **low** severity issue but can lead to subtle bugs or misinterpretation of the data in some APIs.

**Location:** E.g., `BingChatEnabled` in multiple structs including `EnviromentPropertiesDto`, `GenerativeAiFeaturesPropertiesDto`, etc.

**Code Issue:**

```go
BingChatEnabled bool `json:"bingChatEnabled,omitempty"`
```

**Fix:** Consider making boolean fields pointers (`*bool`) to distinguish unset from false:

```go
BingChatEnabled *bool `json:"bingChatEnabled,omitempty"`
```

## ISSUE 5

### Inconsistent Go Struct Tag Formatting

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/dto.go`

**Problem:** Some struct tags are missing a space after the field type and before the struct tag, which is the idiomatic Go style. For example, many lines look like:

```go
X string`json:"x,omitempty"`
```

Rather than:

```go
X string `json:"x,omitempty"`
```

**Impact:** While this does not break compilation or runtime correctness, it is less readable and less idiomatic, and may annoy code reviewers or trigger style linters. Severity: **low**.

**Location:** Check all data structs for missing spaces before struct tags throughout the file.

**Code Issue:**

```go
Id string`json:"id"`
```

**Fix:** Add a space between the field definition and the struct tag in all instances.

```go
Id string `json:"id"`
```

## ISSUE 6

### Missing 'omitempty' on Required JSON Field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/dto.go`

**Problem:** In the structure `EnviromentPropertiesDto`, the field `Description` does not have an `omitempty` in its JSON tag, which means it will always be present in the marshaled JSON, even if set to the empty string. Most other fields use `omitempty`, suggesting this was likely unintentional.

**Impact:** This causes inconsistent API responses and can be confusing. API consumers may expect similar field presence/absence semantics for all optional fields. This is a medium-severity data consistency issue.

**Location:** Struct `EnviromentPropertiesDto`, field `Description`, likely around line 51

**Code Issue:**

```go
Description string `json:"description"`
```

**Fix:** Update this field to use `omitempty` for consistency:

```go
Description string `json:"description,omitempty"`
```

## ISSUE 7

### Field Name and JSON Tag Mismatch

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go`

**Problem:** The Go struct field `PowerAppsComponentFrameworkForCanvasApps` does not match the JSON tag `iscustomcontrolsincanvasappsenabled`, reducing maintainability and readability. The field name should more closely reflect the JSON tag or domain vocabulary for clarity and consistency.

**Impact:** This inconsistency can confuse developers and cause maintenance issues, especially when generating or mapping from API documentation. **Severity:** low.

**Location:**

```go
PowerAppsComponentFrameworkForCanvasApps *bool `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

**Code Issue:**

```go
PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

**Fix:** Rename the field so it matches the domain concept and JSON tag:

```go
IsCustomControlsInCanvasAppsEnabled *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

## ISSUE 8

### Missing JSON Tag on DTO Field

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go`

**Problem:** The struct `linkedEnvironmentIdMetadataDto` has a field `InstanceURL` with no JSON tag. This may cause incorrect mapping if the JSON payload uses a different casing or name for the instance URL, which breaks deserialization/serialization.

**Impact:** May cause JSON parsing issues when the field name does not match exactly what is in the JSON payload (Go's default is to lowercase the struct field name for JSON mapping, but if the server responds with a different case, this will break). **Severity:** medium.

**Location:**

```go
type linkedEnvironmentIdMetadataDto struct {
    InstanceURL string
}
```

**Code Issue:**

```go
InstanceURL string
```

**Fix:** Add a JSON tag reflecting the actual key in the JSON response (e.g. `instanceUrl`). Adjust according to the actual API.

```go
InstanceURL string `json:"instanceUrl"`
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
