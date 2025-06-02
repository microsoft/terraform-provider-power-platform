# Title

Deeply Nested Data Structures

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go

## Problem

The `AnalyticsDataDto` struct includes deeply nested fields, such as the `Environments`, `Status`, and `Sink` structs. While this may work for smaller data sets or scenarios with minimal complexity, deeply nested data structures can make the code harder to read, maintain, and extend over time.

## Impact

Deeply nested data structures can:

- Make the code harder to debug and understand.
- Lead to performance issues when working with large data sets.
- Increase the risk of introducing bugs when these structs are modified or extended.

Severity: **medium**

## Location

File: `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go`

```go
type AnalyticsDataDto struct {
    ID               string           `json:"id"`
    Source           string           `json:"source"`
    Environments     []EnvironmentDto `json:"environments"`
    Status           []StatusDto      `json:"status"`
    Sink             SinkDto          `json:"sink"`
    PackageName      string           `json:"packageName"`
    Scenarios        []string         `json:"scenarios"`
    ResourceProvider string           `json:"resourceProvider"`
    AiType           string           `json:"aiType"`
}

```

## Fix

Consider refactoring this struct to:

1. Use helper methods to process deeply nested data.
2. Decompose the struct into smaller, more focused modules or interfaces that manage their respective functionality.

Example Alteration:

```go
func (dto AnalyticsDataDto) GetEnvironmentIds() []string {
    var environmentIds []string
    for _, env := range dto.Environments {
        environmentIds = append(environmentIds, env.EnvironmentId)
    }
    return environmentIds
}
```