# Struct and receiver naming does not follow Go conventions

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The `client` struct and its receiver variable are lowercase and generic, conflicting with Go naming conventions. According to Go best practices:
- Type names should be upper camel case (exported/public) or at least capitalized internally to distinguish types from variables.
- Receiver names should be short, typically the first letter or abbreviation of the type (`c` for `Client`).

The current conventions may create readability and maintainability issues, especially as the codebase grows and `client` could clash with other variables or imports.

## Impact

Severity: Medium. Hinders code readability/maintainability and risks future naming conflicts.

## Location

Type and method declarations:

## Code Issue

```go
type client struct {
	Api *api.Client
}

func (client *client) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
   // ...
}

func (client *client) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
   // ...
}
```

## Fix

Rename the struct to `Client` and update receivers from `client` to `c`:

```go
type Client struct {
	Api *api.Client
}

func (c *Client) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
   // ...
}

func (c *Client) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
   // ...
}
```
