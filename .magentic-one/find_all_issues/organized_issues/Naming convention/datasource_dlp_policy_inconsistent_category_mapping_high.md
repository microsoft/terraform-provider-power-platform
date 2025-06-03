# Inconsistent Category Mapping for Connector Groups

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

In the iteration over `policies` in the `Read` function, there is a possible mistake in how the conversion helpers are called:

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
```

Naming implies that "BusinessGeneralConnectors" should relate to "General" and "NonBusinessConfidentialConnectors" to "Confidential". However, the code uses the opposite mapping, potentially leading to incorrect data population.

## Impact

This could result in connectors being assigned to the wrong data group in the Terraform state, leading to incorrect resource configuration and potentially violating data loss prevention logic. Severity: **high**.

## Location

Function: `Read`
Lines where `convertToAttrValueConnectorsGroup` is called.

## Code Issue

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
```

## Fix

Swap the category strings where the conversion functions are called, so the variable names match the values:

```go
policyModel.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
policyModel.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
```

**Explanation:**  
This aligns the data being set with the correct expected category as described by the field names, ensuring the right connectors are mapped.
