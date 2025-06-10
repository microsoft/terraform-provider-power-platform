# DTO Issues - Input Validation

This document contains all DTO-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Title

Data Consistency: Use of String for Date Fields Instead of Time Type

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/dto.go

## Problem

Fields such as `CreatedOn` and `LastModifiedOn` in `BillingPolicyDto` are defined as `string` instead of using Go's `time.Time` (with the standard `encoding/json` support via the `time` package). Using `string` for date/time values can lead to inconsistent date formats, lack of validation, and more error-prone handling of these fields.

## Impact

This can introduce bugs in date/time manipulation, cause inconsistency in how dates are serialized or deserialized, and make validation harder. Severity is **medium** as it affects data consistency and type safety, though it might also be API-driven.

## Location

`BillingPolicyDto` struct:

## Code Issue

```go
type BillingPolicyDto struct {
	// ...
	CreatedOn         string               `json:"createdOn"`
	CreatedBy         PrincipalDto         `json:"createdBy"`
	LastModifiedOn    string               `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto         `json:"lastModifiedBy"`
}
```

## Fix

Use the `time.Time` type for date fields, and ensure that JSON marshaling/unmarshaling is handled correctly (e.g., via RFC3339 or whatever format your API expects).

```go
import "time"

type BillingPolicyDto struct {
	// ...
	CreatedOn         time.Time            `json:"createdOn"`
	CreatedBy         PrincipalDto         `json:"createdBy"`
	LastModifiedOn    time.Time            `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto         `json:"lastModifiedBy"`
}
```

If custom formatting is required, provide the appropriate (un)marshal methods. Only use `string` if the upstream API dictates it *and* there is no ability to change serialization behavior.


## ISSUE 2

# Title

Async, Validate, and Environment ID DTOs Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The types `asyncSolutionPullResponseDto`, `validateSolutionImportResponseDto`, `validateSolutionImportResponseSolutionOperationResultDto`, `environmentIdDto`, `environmentIdPropertiesDto`, and `linkedEnvironmentIdMetadataDto` are all unexported, likely limiting their usability for consumers of the package with type embedding, function parameters, or API responses. DTOs should generally be exported for broad usability.

## Impact

Not exporting these types limits code extensibility, clarity, and reuse. Severity: medium.

## Location

Lines 100–127

## Code Issue

```go
type asyncSolutionPullResponseDto struct { ... }
type validateSolutionImportResponseDto struct { ... }
type validateSolutionImportResponseSolutionOperationResultDto struct { ... }
type environmentIdDto struct { ... }
type environmentIdPropertiesDto struct { ... }
type linkedEnvironmentIdMetadataDto struct { ... }
```

## Fix

Capitalize each type to make them exported:

```go
type AsyncSolutionPullResponseDto struct { ... }
type ValidateSolutionImportResponseDto struct { ... }
type ValidateSolutionImportResponseSolutionOperationResultDto struct { ... }
type EnvironmentIdDto struct { ... }
type EnvironmentIdPropertiesDto struct { ... }
type LinkedEnvironmentIdMetadataDto struct { ... }
```


## ISSUE 3

# Title

Should Use time.Time For Date/Time Fields

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

Several fields representing date/time values are typed as `string` (e.g., `CreatedTime`, `ModifiedTime`, `InstallTime`, `CreatedOn`, `CompletedOn`). Using `string` does not provide compile-time type safety and leaves the responsibility of handling and validating date/time formats to the consumers. Using `time.Time` provides better semantic expressiveness, parsing, and validation benefits. The encoding/json package supports custom (un)marshalers for time.Time types.

## Impact

Risk of invalid date/time string values, lack of compile-time checks, possible bugs when processing or comparing date/times. Severity: medium.

## Location

Lines 22–24, 103–105

## Code Issue

```go
CreatedTime   string `json:"createdon"`
ModifiedTime  string `json:"modifiedon"`
InstallTime   string `json:"installedon"`
...
CreatedOn        string `json:"createdon"`
CompletedOn      string `json:"completedon"`
```

## Fix

Use `time.Time` as the type for fields representing date/time values. If non-standard formats are present, implement `UnmarshalJSON` methods. Example:

```go
import "time"

CreatedTime   time.Time `json:"createdon"`
ModifiedTime  time.Time `json:"modifiedon"`
InstallTime   time.Time `json:"installedon"`
```

When unmarshalling, ensure that the time format matches your JSON. Example custom unmarshal logic might be needed if format differs from RFC3339.


## ISSUE 4

# Use of `string` for Times Instead of `time.Time`

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/dto.go

## Problem

The fields `CreatedTime`, `LastModifiedTime`, and `LastPublishTime` in `powerAppPropertiesBapiDto` are typed as `string`. In Go, it's best practice to use the `time.Time` type for date/time data to ensure proper parsing, validation, and time operations, instead of using raw strings.

## Impact

Severity: Medium

This may lead to data inconsistency, missed validation (like parsing), inability to use Go's time utilities, and bugs related to time formatting, localization, or comparison.

## Location

```go
type powerAppPropertiesBapiDto struct {
    // ...
    CreatedTime      string                 `json:"createdTime"`
    LastModifiedTime string                 `json:"lastModifiedTime"`
    LastPublishTime  string                 `json:"lastPublishTime"`
    // ...
}
```

## Code Issue

```go
CreatedTime      string                 `json:"createdTime"`
LastModifiedTime string                 `json:"lastModifiedTime"`
LastPublishTime  string                 `json:"lastPublishTime"`
```

## Fix

Use `time.Time` instead of `string` and, if necessary, provide custom JSON unmarshal/marshal logic if the time format is not RFC3339.

```go
import "time"

type powerAppPropertiesBapiDto struct {
    // ...
    CreatedTime      time.Time              `json:"createdTime"`
    LastModifiedTime time.Time              `json:"lastModifiedTime"`
    LastPublishTime  time.Time              `json:"lastPublishTime"`
    // ...
}
```

If a custom date/time format is used, implement `UnmarshalJSON` and `MarshalJSON` methods for proper parsing.


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
