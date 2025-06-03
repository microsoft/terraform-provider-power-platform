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
