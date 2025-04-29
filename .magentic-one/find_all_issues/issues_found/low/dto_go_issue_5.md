# Title

Unnecessary Boilerplate Code

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go

## Problem

Struct declarations like `AnalyticsDataResponse`, `AnalyticsDataDto`, and others include repetitive and generic fields without leveraging common Go techniques such as embedding other types or utilizing interfaces. This increases the codebase size unnecessarily and may introduce duplication.

## Impact

Unnecessary boilerplate:

- Makes the code harder to refactor and maintain.
- Wastes developer resources when implementing future changes.
- Reduces readability and clarity of the codebase.

Severity: **low**

## Location

File: `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go`

```go
type AnalyticsDataResponse struct {
    Value []AnalyticsDataDto `json:"value"`
}

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

Reduce boilerplate by embedding shared struct fields where applicable. For example:

```go
type BaseDto struct {
    ID   string `json:"id"`
    Sink SinkDto `json:"sink"`
}

type AnalyticsDataDto struct {
    BaseDto
    Source           string           `json:"source"`
    Environments     []EnvironmentDto `json:"environments"`
    Status           []StatusDto      `json:"status"`
    PackageName      string           `json:"packageName"`
    Scenarios        []string         `json:"scenarios"`
    ResourceProvider string           `json:"resourceProvider"`
    AiType           string           `json:"aiType"`
}
```