# Title

Potential JSON Tag Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go

## Problem

In the `SinkDto` struct, some fields use the `omitempty` directive in their JSON tags (e.g., `SubscriptionId`, `ResourceGroupName`), while others do not. This inconsistency can lead to unclear behavior when the struct is serialized to JSON, especially in cases where fields with zero or empty values should not be included.

## Impact

Inconsistent JSON behavior may:

- Cause issues during integration with external systems or APIs.
- Lead to unexpected behavior or bugs when certain fields are omitted unintentionally.
- Increase confusion for developers working with the serialization logic.

Severity: **high**

## Location

File: `/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go`

```go
type SinkDto struct {
    ID                string `json:"id"`
    Type              string `json:"type"`
    SubscriptionId    string `json:"subscriptionId,omitempty"`
    ResourceGroupName string `json:"resourceGroupName,omitempty"`
    ResourceName      string `json:"resourceName"`
    Key               string `json:"key"`
}
```

## Fix

Ensure that all fields either consistently use `omitempty` if they can be omitted when empty, or remove `omitempty` to enforce inclusion in the JSON output. For example:

```go
type SinkDto struct {
    ID                string `json:"id,omitempty"`
    Type              string `json:"type,omitempty"`
    SubscriptionId    string `json:"subscriptionId,omitempty"`
    ResourceGroupName string `json:"resourceGroupName,omitempty"`
    ResourceName      string `json:"resourceName,omitempty"`
    Key               string `json:"key,omitempty"`
}

// Or if omitting is not required:
type SinkDto struct {
    ID                string `json:"id"`
    Type              string `json:"type"`
    SubscriptionId    string `json:"subscriptionId"`
    ResourceGroupName string `json:"resourceGroupName"`
    ResourceName      string `json:"resourceName"`
    Key               string `json:"key"`
}
```