# Type Assertion Without Validation in `Configure` Method

## Path

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

In the `Configure` method, the code attempts a type assertion without validating if the type assertion was successful. If the `ProviderData` provided does not match a `*api.ProviderClient` type, the code will encounter a runtime panic instead of failing gracefully.

## Impact

- **Severity: Critical**
- Runtime panic on invalid type assertion can crash the program and affect the stability of the entire Terraform provider.
- Poor error handling may lead to unexpected client behavior and decrease trust in the software.

## Location

```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}
	clientApi := req.ProviderData.(*api.ProviderClient).Api // Issue here
	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.LanguagesClient = newLanguagesClient(clientApi)
}
```

## Fix

The type assertion should be validated before attempting to use the derived type. This avoids runtime panics and ensures robust error handling.

### Revised Code
```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient) // Validate type assertion
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Provider Data Type",
			fmt.Sprintf("Expected *api.ProviderClient type but got: %T. If this issue persists, please report it.", req.ProviderData),
		)
		return
	}

	clientApi := providerClient.Api
	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Provider client API object is nil. Please report this issue to the provider developers.",
		)
		return
	}

	d.LanguagesClient = newLanguagesClient(clientApi)
}
```

### Explanation
1. Added a check using `ok` to ensure the type assertion is valid.
2. Added meaningful error messages to assist users in identifying issues.
3. Prevented runtime panics, ensuring better reliability.
