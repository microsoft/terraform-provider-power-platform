# Title

Missing Comma in `jsonencode` Blocks

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

This issue occurs in the `connection_parameters` property in both test cases. A missing comma separates key-value pairs, leading to invalid JSON syntax.

#### Violating Code:

```go
connection_parameters = jsonencode({
	"azureOpenAIResourceName" : "aaa",
	"azureOpenAIApiKey" : "bbb"
	"azureSearchEndpointUrl" : "ccc",
	"azureSearchApiKey" : "ddd"
})
```

## Impact

- Severity: Critical
- Results in runtime parsing errors. Deployment with this flawed code will fail when creating the connection resource.

## Location

Two occurrences found:

1. Part of `TestAccConnectionsResource_Validate_Create` function at line 22.
2. Part of `TestUnitConnectionsResource_Validate_Create` function at line 86.

## Fix

Correct the JSON syntax by adding commas between each key-value pair:

```go
connection_parameters = jsonencode({
	"azureOpenAIResourceName" : "aaa",
	"azureOpenAIApiKey" : "bbb",
	"azureSearchEndpointUrl" : "ccc",
	"azureSearchApiKey" : "ddd"
})
```