# Title

Missing Commas in JSON Object Literals in Terraform Config

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

In both tests under the `Config` fields, the `connection_parameters` attribute is given a JSON object that is missing a comma between `"azureOpenAIApiKey" : "bbb"` and `"azureSearchEndpointUrl" : "ccc",`. The current code:

```json
"azureOpenAIResourceName" : "aaa",
"azureOpenAIApiKey" : "bbb"
"azureSearchEndpointUrl" : "ccc",
"azureSearchApiKey" : "ddd"
```

should have a comma after `"azureOpenAIApiKey" : "bbb"`.

## Impact

Malformed JSON will cause a runtime error when the Terraform provider tries to parse the string, leading to failing tests that do not actually test the functionality or the error handling of the code being tested. **Severity: High**

## Location

Both in the `Config` field of each test, in the definition of `connection_parameters`.

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

Add a comma after `"azureOpenAIApiKey" : "bbb"` in both test cases.

```go
connection_parameters = jsonencode({
	"azureOpenAIResourceName" : "aaa",
	"azureOpenAIApiKey" : "bbb",
	"azureSearchEndpointUrl" : "ccc",
	"azureSearchApiKey" : "ddd"
})
```
