# Title

Method Receiver Should Be Pointer to Avoid Unintended Copy

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/client.go

## Problem

The method `GetSolutionCheckerRules` uses a non-pointer receiver (`c *client`). This is correct. However, if `client` is used without pointer semantics (because the constructor returns a value), the receiver and method call will work on a copy rather than the original, potentially leading to bugs if the method ever modifies state. The current constructor returns a value, so clients may end up using non-pointer semantics by mistake. This is a code structure issue.

## Impact

Potential unintended copies and inconsistent usage. If the struct gets fields that need mutation or synchronization, could cause bugs. Severity: low to medium.

## Location

- Constructor and receiver usage for methods on `client`

## Code Issue

```go
func newSolutionCheckerRulesClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

func (c *client) GetSolutionCheckerRules(ctx context.Context, environmentId string) ([]ruleDto, error) {
	//...
}
```

## Fix

Return a pointer from the constructor, and ensure all usages are as pointer. This is considered good Go client API practice.

```go
func NewSolutionCheckerRulesClient(apiClient *api.Client) *Client {
	return &Client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

func (c *Client) GetSolutionCheckerRules(ctx context.Context, environmentId string) ([]ruleDto, error) {
	//...
}
```
