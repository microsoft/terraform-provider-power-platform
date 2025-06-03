# Reliance on External Files for Testing

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

The tests rely on the presence of JSON files under `tests/datasource/Web_Apis_WhoAmI/`. If these files are missing or the relative path changes, the test will crash or silently pass with incorrect data (via empty .String()).

## Impact

Medium severity for test reliability and robustness. This introduces unnecessary fragility.

## Location

```go
httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json")
```

## Fix

Where possible, inline small JSON strings as test data, or add setup/teardown routines to clearly fail if a dependency file is missing. At minimum, ensure error handling is robust, and add comments to clarify the dependency.

