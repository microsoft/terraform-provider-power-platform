# Title

Incorrect usage of field names in ValidateConfig (BusinessGeneralConnectors, NonBusinessConfidentialConnectors)

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

The function `ValidateConfig` references struct fields `BusinessGeneralConnectors`, `NonBusinessConfidentialConnectors`, and `BlockedConnectors`. However, the model schema, as defined in `Schema`, uses the field names `business_connectors`, `non_business_connectors`, and `blocked_connectors`. The use of incorrect field names will result in runtime errors or unexpected behavior as the fields may not be populated or marshaled from Terraform input properly.

## Impact

Severity: High

This issue impacts the correctness of validation logic and may prevent configuration validation on Terraform plans. Any validation using these incorrect names will not operate as expected and may result in misconfigured DLP policies.

## Location

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessGeneralConnectors)
...
conn, err = getConnectorGroup(ctx, config.NonBusinessConfidentialConnectors)
...
conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
...
```

## Code Issue

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessGeneralConnectors)
if err != nil {
	resp.Diagnostics.AddError("BusinessGeneralConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.NonBusinessConfidentialConnectors)
if err != nil {
	resp.Diagnostics.AddError("NonBusinessConfidentialConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
if err != nil {
	resp.Diagnostics.AddError("BlockedConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)
```

## Fix

Replace `BusinessGeneralConnectors` with `BusinessConnectors`, `NonBusinessConfidentialConnectors` with `NonBusinessConnectors` to align with the schema and struct field names.

```go
var connectors []dlpConnectorModelDto
conn, err := getConnectorGroup(ctx, config.BusinessConnectors)
if err != nil {
	resp.Diagnostics.AddError("BusinessConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.NonBusinessConnectors)
if err != nil {
	resp.Diagnostics.AddError("NonBusinessConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)

conn, err = getConnectorGroup(ctx, config.BlockedConnectors)
if err != nil {
	resp.Diagnostics.AddError("BlockedConnectors validation error", err.Error())
}
connectors = append(connectors, conn.Connectors...)
```
---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_dlp_policy_incorrect_field_name_high.md
