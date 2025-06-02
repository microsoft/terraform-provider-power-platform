# Insufficient Negative Test Coverage

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

The tests focus only on positive scenarios where the API call and file read succeed. There are no tests for negative scenarios, such as HTTP errors, missing/mock files, or invalid response formats.

## Impact

Medium severity as this reduces test coverage and could allow undetected bugs related to error handling, input validation, or unexpected failure cases.

## Location

No explicit code location, but the absence of negative testing is an issue.

## Fix

Add additional test cases that deliberately trigger failure modes, such as:

- Mock responder returns 500 or invalid JSON.
- Missing or corrupt JSON test file.
- API returns unexpected data format.

Add assertions to ensure the error is handled gracefully and correctly reported to the user.

