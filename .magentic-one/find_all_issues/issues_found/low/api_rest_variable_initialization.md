# Title

Unnecessary Initialization of Variables in `SendOperation`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In the `SendOperation` function, variables like `body` and `headers` are unnecessarily initialized or set with additional checks. For example:
- `body` is declared as `*string` but is left `nil` unless a value is explicitly assigned later.
- `headers` is defined only if `operation.Headers` has a length greater than zero.

This kind of variable initialization creates redundant code that can complicate readability and introduces unnecessary decision branches.

## Impact

- Reduces code readability: Unnecessary logic adds unneeded complexity to the function.
- Wastes developer time in debugging or understanding the flow.
- Can slightly degrade performance if multiple tests and assignments are made unnecessarily.

Severity: **Low**

## Location

Found in `SendOperation`.

## Code Issue

```go
var body *string
var headers map[string]string
if operation.Body.ValueStringPointer() != nil {
	b := operation.Body.ValueString()
	body = &b
}
if len(operation.Headers) > 0 {
	headers = make(map[string]string)
	for _, h := range operation.Headers {
		headers[h.Name.ValueString()] = h.Value.ValueString()
	}
}
```

## Fix

Streamline the initialization process to avoid unnecessary declarations or logic. Use more concise conditional assignments.

```go
body := func() *string {
	if operation.Body.ValueStringPointer() != nil {
		b := operation.Body.ValueString()
		return &b
	}
	return nil
}()

headers := func() map[string]string {
	if len(operation.Headers) > 0 {
		h := make(map[string]string)
		for _, header := range operation.Headers {
			h[header.Name.ValueString()] = header.Value.ValueString()
		}
		return h
	}
	return nil
}()
```

This improves readability and ensures variables are initialized only when required.