# Title

Testing: No coverage or tests for edge-cases and error branches

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

There are no indications of unit tests or error-branch validation for:
- Type assertion edge-cases (failed or unexpected API response types)
- Error branches for each conditional return
- API status codes outside expected cases

## Impact

**Severity: Medium**

- Increases risk of regression or shipping unnoticed bugs when refactoring.
- Undocumented behavior for how the code handles API changes or malformed responses.
- Makes automated refactoring or security-tuning riskier.

## Location

N/A (absence of test code, but logic such as this needs proper test coverage):

```go
if response["@odata.context"] == nil { ... }
if !ok { return nil, errors.New(...) }
if response.HttpResponse.StatusCode == http.StatusPreconditionFailed { ... }
```

## Code Issue

N/A

## Fix

Add focused unit tests for each function using a table-driven approach that covers:

- HTTP/JSON happy-path and typical failures
- Each error return, including missing fields and failed type assertions
- Handling for unexpected status codes, body shapes, or nils

Consider using a fake/mock for `client.Api` and for responses to drive branch and edge case coverage.

---

File:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/api_data_record_no_edgecase_tests_medium.md`
