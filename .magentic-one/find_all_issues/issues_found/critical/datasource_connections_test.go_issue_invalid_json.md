# Title

Improperly Escaped String in `connection_parameters` Key Block

##

`/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go`

## Problem

In the `connection_parameters` JSON block in `TestAccConnectionsDataSource_Validate_Read`, the JSON is not properly escaped before being passed as a string within the Terraform configuration. This is evidenced by missing commas between keys, which invalidates the JSON string.

For instance: 
```hcl
"azureOpenAIApiKey" : "bbb"
"azureSearchEndpointUrl" : "ccc",
```
lacks a comma between `"bbb"` and `"ccc"`. This will result in runtime errors.

## Impact

This issue will result in the Terraform configuration failing to parse properly, leading to test execution errors. Severity: **Critical**, as it directly impacts the ability for users to run the tests effectively.

## Location

- Test function: `TestAccConnectionsDataSource_Validate_Read`
- Line containing `connection_parameters` JSON encoding.

## Code Issue

```go
connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "aaa",
    "azureOpenAIApiKey" : "bbb"
    "azureSearchEndpointUrl" : "ccc",
    "azureSearchApiKey" : "ddd"
})
```

## Fix

Correct the JSON formatting by ensuring commas are used between all key-value pairs:

```go
connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "aaa",
    "azureOpenAIApiKey" : "bbb",
    "azureSearchEndpointUrl" : "ccc",
    "azureSearchApiKey" : "ddd"
})
```

This change ensures the JSON block is syntactically valid. Any downstream operations relying on the validity of this JSON data will now function correctly.
