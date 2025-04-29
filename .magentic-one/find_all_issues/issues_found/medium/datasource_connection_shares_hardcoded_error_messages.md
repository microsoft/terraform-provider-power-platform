# Title

Hardcoded Diagnostic Error Messages without Localization

##

`/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go`

## Problem

The diagnostic error messages in the `Configure` method of the `SharesDataSource` are hardcoded, including `"Unexpected ProviderData Type"` and specifics of the expected and received data types (`fmt.Sprintf("Expected *api.ProviderClient, got: %T", req.ProviderData)`). This practice does not adhere to localization standards or extensible error handling methods.

## Impact

Hardcoded error messages reduce the maintainability of the codebase, particularly when the application needs to support multiple languages or dynamically change messages based on user context. Severity: **medium**.

## Location

Code File Path: `/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go`, within the `Configure` function:

## Code Issue

```go
if !ok {
	resp.Diagnostics.AddError(
		"Unexpected ProviderData Type",
		fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
	return
}
```

## Fix

Define the error messages in a central localization or error handling file/library, and use them in the code as demonstrated below:

```go
const (
	UnexpectedProviderDataErrorTitle = "Unexpected Provider Data Type"
	UnexpectedProviderDataErrorDetails = "Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers."
)

if !ok {
	resp.Diagnostics.AddError(
		UnexpectedProviderDataErrorTitle,
		fmt.Sprintf(UnexpectedProviderDataErrorDetails, req.ProviderData),
	)
	return
}
```

Explanation:
- Centralize error messages in language-agnostic constants or external resource files to allow easy localization or updates.
- Use placeholders (`%T`) to dynamically insert runtime values into error messages.

---
