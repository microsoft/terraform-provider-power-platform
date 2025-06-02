# Test Contains Typo in JSON Syntax in HCL Inline String

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares_test.go

## Problem

Within the large HCL config in the acceptance test, the `connection_parameters` JSON object includes a missing comma between `"azureOpenAIApiKey" : "bbb"` and `"azureSearchEndpointUrl" : "ccc",`. This results in invalid JSON, which may cause errors when the test is run and the configuration is parsed, despite Go's string handling not flagging this directly.

## Impact

Severity: **High**  
This is likely to cause test failures due to invalid configuration passed to Terraform. It also demonstrates the risk of embedding large, non-validated blobs in Go code without clear error reporting.

## Location

```go
connection_parameters = jsonencode({
  "azureOpenAIResourceName" : "aaa",
  "azureOpenAIApiKey" : "bbb"
  "azureSearchEndpointUrl" : "ccc",
  "azureSearchApiKey" : "ddd"
})
```

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

Insert a comma after `"azureOpenAIApiKey" : "bbb"` to ensure the JSON is valid:

```go
connection_parameters = jsonencode({
  "azureOpenAIResourceName" : "aaa",
  "azureOpenAIApiKey" : "bbb",
  "azureSearchEndpointUrl" : "ccc",
  "azureSearchApiKey" : "ddd"
})
```
