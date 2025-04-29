# Title
Hardcoding of Metadata Context URLs

## Problem
Several URLs, such as `https://api.bap.microsoft.com/providers/...` and `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/...`, are hardcoded within the test file.

## Impact
Hardcoding such URLs makes the code less maintainable and dependent on specific runtime conditions. If these URLs change in the future, the tests will break, requiring updates across the file. This is a **high severity** issue when working with dynamic configurations and services.

## Location
In `resource_rest_test.go`, lines defining `beforeUpdateRegex` and `afterUpdateRegex` and other locations where URLs are hardcoded.

## Code Issue
```go
beforeUpdateRegex := `^	{"@odata\.context":"https:\/\/org[0-9a-fA-F]{8}\.crm\.dynamics\.com\/api\/data\/v9\.2\/\$metadata#accounts\(name,accountid\)\/\$entity","@odata\.etag":"W\/\"[0-9]{7}\"","name":"powerplatform_rest","accountid":"00000000-0000-0000-0000-000000000001"}`
afterUpdateRegex := `^	{"@odata\.context":"https:\/\/org[0-9a-fA-F]{8}\.crm\.dynamics\.com\/api\/data\/v9\.2\/\$metadata#accounts\(name,accountid\)\/\$entity","@odata\.etag":"W\/\"[0-9]{7}\"","name":"powerplatform_rest_change","accountid":"00000000-0000-0000-0000-000000000001"}`
```

## Fix
Define constants in a configuration or mock setup file, ensuring URLs can be reused without being hardcoded directly into your tests. This also ensures easy updates when URLs change.

```go
// Mock constants for URLs
const baseURL = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com"
const metadataContextURL = fmt.Sprintf("%s/api/data/v9.2/$metadata#accounts", baseURL)

// Modify variable setup
beforeUpdateRegex := fmt.Sprintf(`{"@odata\.context":"%s...")
```