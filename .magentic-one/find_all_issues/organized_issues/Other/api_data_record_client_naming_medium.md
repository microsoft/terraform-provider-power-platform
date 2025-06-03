# Title

Naming convention: Struct type `client` does not follow Go best practice (should be `Client`)

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

The struct type representing the API client is defined with a lowercase `client`, which restricts its visibility to the package and does not follow Go's widely accepted convention for exported types. Although it may be intended for internal use, this struct is returned and referenced in public API signatures (for methods and factory).

## Impact

**Severity: Medium**

- Reduces team readability, as "client" types are normally capitalized for clarity and to distinguish as type, not variable.
- If there is a plan to expand package public API in the future, not aligning to conventions causes confusion and technical debt.
- Makes refactoring harder if the type needs to eventually be exported.

## Location

```go
type client struct {
	Api *api.Client
}
```

## Code Issue

```go
type client struct {
	Api *api.Client
}
```

## Fix

Capitalize the struct name for clarity and to meet Go conventions (unless strictly required for encapsulation):

```go
type Client struct {
	Api *api.Client
}
```

You will also need to update all method receivers and usages:

```go
func (c *Client) MethodName(...) ...
```

---

Save as:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_data_record_client_naming_medium.md`
