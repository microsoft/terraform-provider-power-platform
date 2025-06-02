# Issue 2

Misformatted JSON in Resource Configuration String

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go

## Problem

In the configuration string for `"powerplatform_connection" "azure_openai_connection"`, the JSON object within `jsonencode({ ... })` is missing a comma after `"azureOpenAIApiKey" : "bbb"`. This is a syntax error and may cause test failures or unexpected test behavior.

## Impact

- **Severity:** High
- Will cause Terraform or the test runner to fail parsing the configuration.
- The test does not validate the scenario it is supposed to, causing misleading test results or false failures.
- Test cannot run until syntax has been corrected.

## Location

```go
connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "aaa",
    "azureOpenAIApiKey" : "bbb"
    "azureSearchEndpointUrl" : "ccc",
    "azureSearchApiKey" : "ddd"
})
```

## Fix

Insert a comma after `"azureOpenAIApiKey" : "bbb"` so that the JSON encode input is valid.

```go
connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "aaa",
    "azureOpenAIApiKey" : "bbb",
    "azureSearchEndpointUrl" : "ccc",
    "azureSearchApiKey" : "ddd"
})
```
This ensures the configuration is syntactically correct and can be processed properly by Terraform and the test infrastructure.
