# Title

Incorrect HTTP status code validation for policy linking and unlinking operations.

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The current implementation for validating HTTP status codes in the response uses a hardcoded list (`http.StatusAccepted`) of acceptable status codes. However, it does not account for additional potential success codes (e.g., `http.StatusOK`) or handle potential edge cases where the status code may vary depending on the API implementation.

## Impact

This problem can lead to false negatives in response validation, causing unnecessary retry attempts or exceptions, thereby potentially lowering system reliability and user trust. Severity of this issue is **medium**.

## Location

- `LinkEnterprisePolicy`: `apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)`
- `UnLinkEnterprisePolicy`: `apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)`

## Code Issue

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)
```

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)
```

## Fix

Expand the list of accepted HTTP status codes to encompass all possible success statuses. For instance, include `http.StatusOK` and any other codes defined in the API documentation as indicative of success. 

```go
apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprisePolicyDto, []int{http.StatusAccepted, http.StatusOK}, nil)
```

If additional success logic for non-standard status codes is required, implement a utility function to dynamically validate response codes based on API expectations. Improve comments and documentation for clarity. For example:

```go
func isValidHttpStatus(statusCode int) bool {
	acceptedStatusCodes := []int{http.StatusAccepted, http.StatusOK}
	for _, validCode := range acceptedStatusCodes {
		if statusCode == validCode {
			return true
		}
	}
	return false
}
``` 

Alter the code to use this function:

```go
if !isValidHttpStatus(apiResponse.HttpResponse.StatusCode) {
	return fmt.Errorf("Invalid HTTP response status %d", apiResponse.HttpResponse.StatusCode)
}
```

This fix improves resilience against API-specific edge cases and strengthens code quality standards.
