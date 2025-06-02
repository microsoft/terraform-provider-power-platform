# Title

No Error Handling for Mock Responder File Loading

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

In the registered responder, the use of `httpmock.File(...).String()` assumes the file will always be present and can be read. If the file is missing or unreadable, the tests may silently fail with confusing errors or panics.

## Impact

Medium. This impacts test reliability. If test data is missing or inaccessible, the test should fail gracefully with descriptive errors so the issue is apparent.

## Location

Lines 17–19.

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant_capacity.json").String()), nil
```

## Fix

Explicitly check for file loading errors before using the contents. If your test helper does not return errors, consider checking for file existence first or handle the panic explicitly.

```go
contentBytes, err := os.ReadFile(mockCapacityFile)
if err != nil {
	t.Fatalf("failed to read mock response file: %v", err)
}
return httpmock.NewStringResponse(http.StatusOK, string(contentBytes)), nil
```
