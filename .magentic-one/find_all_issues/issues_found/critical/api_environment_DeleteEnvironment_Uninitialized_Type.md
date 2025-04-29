# Title

Incorrect reference to uninitialized structure in `DeleteEnvironment`

##

`/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

## Problem

In the `DeleteEnvironment` function, the `environmentDelete` structure is used for sending the DELETE request with specific message and code. However, the structure name `enironmentDeleteDto` is misspelled and hasn't been declared or initialized anywhere in the code, resulting in a compilation error.

## Impact

This issue will prevent the code from compiling successfully, rendering the implementation of the `DeleteEnvironment` function non-functional. Severity: Critical.

## Location

The problem occurs in the following method:

```go
func (client *Client) DeleteEnvironment(ctx context.Context, environmentId string) error
```

## Code Issue

```go
	environmentDelete := enironmentDeleteDto{
		Code:    "7", // Application.
		Message: "Deleted using Power Platform Terraform Provider",
	}
```

## Fix

Correct the typo in the structure name from `enironmentDeleteDto` to the intended type name.

```go
	environmentDelete := environmentDeleteDto{ // Corrected name
		Code:    "7", // Application.
		Message: "Deleted using Power Platform Terraform Provider",
	}
```

Ensure that `environmentDeleteDto` is properly defined and imported, or explicitly declared in the file to avoid compilation errors. For example:

```go
type environmentDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
```