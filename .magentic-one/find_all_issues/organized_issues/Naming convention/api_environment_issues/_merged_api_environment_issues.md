# Api Environment Issues - Merged Issues

## ISSUE 1

# Issue 6: Variable Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go

## Problem

Some receiver names use the full word `client` (`func (client *client)`), which is against Go idiomatic naming—receivers should typically be a single letter (e.g., `c *client`).

## Impact

**Low** severity. Inconsistent receiver naming can reduce code readability and goes against Go best practices.

## Location

Every method defined on `client` struct:

## Code Issue

```go
func (client *client) MethodName(...) {...}
```

## Fix

Refactor all receiver names from `client` to `c`:

```go
func (c *client) MethodName(...) {...}
```


---

## ISSUE 2

# Inconsistent Naming of Struct Fields in client

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go

## Problem

The `client` struct uses mixed naming conventions for its fields (e.g., `Api` vs. `RuleSetApi` vs. `TenantApi`). In Go, consistent naming (all abbreviations either cased as ID/URL/API or not) enhances code readability and maintainability.

## Impact

This reduces code clarity, may lead to confusion, and makes refactoring and code review more difficult.

**Severity:** Low

## Location

```go
type client struct {
	Api        *api.Client
	TenantApi  tenant.Client
	RuleSetApi environment_group_rule_set.Client
}
```

## Fix

Choose a consistent convention for abbreviations (preferably "API" as all-caps, conforming to Go idioms for initialisms).

```go
type client struct {
	API        *api.Client
	TenantAPI  tenant.Client
	RuleSetAPI environment_group_rule_set.Client
}
```

And update all usages accordingly in the whole codebase


---

## ISSUE 3

# Naming: Struct Type Name `Client` Is Too Generic

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

The struct type `Client` is overly generic. This can create confusion and naming collisions, especially within a large codebase or when referencing across packages. Good practice is to use a more context-specific name like `EnvironmentGroupRuleSetClient`.

## Impact

Negatively affects code readability, maintainability, and discoverability. Makes it hard to search for relevant client structs in the codebase. Severity: Low to Medium.

## Location

```go
type Client struct {
    Api       *api.Client
    TenantApi tenant.Client
}
```

## Code Issue

```go
type Client struct {
    Api       *api.Client
    TenantApi tenant.Client
}
```

## Fix

Rename the struct to a more meaningfully scoped name, such as `EnvironmentGroupRuleSetClient`.

```go
type EnvironmentGroupRuleSetClient struct {
    Api       *api.Client
    TenantApi tenant.Client
}

// Update constructor accordingly
func NewEnvironmentGroupRuleSetClient(apiClient *api.Client, tenantClient tenant.Client) EnvironmentGroupRuleSetClient {
    return EnvironmentGroupRuleSetClient{
        Api:       apiClient,
        TenantApi: tenantClient,
    }
}
```

---

This issue will be saved in:
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_environment_group_rule_set_naming_low.md


---

## ISSUE 4

# Naming: Parameter Name TenantApi Uses Inconsistent Casing

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

The struct field and constructor parameter are named `TenantApi` (PascalCase) while Go style suggests using initialisms (such as API) in all-caps, making `TenantAPI` preferable. Consistency in naming makes code easier to read and maintain.

## Impact

Affects codebase conformity to Go naming conventions, potentially impacting readability and onboarding for developers familiar with Go standards. Severity: Low.

## Location

```go
type Client struct {
    Api       *api.Client
    TenantApi tenant.Client
}
...
func NewEnvironmentGroupRuleSetClient(apiClient *api.Client, tenantClient tenant.Client) Client {
    return Client{
        Api:       apiClient,
        TenantApi: tenantClient,
    }
}
```

## Code Issue

```go
    TenantApi tenant.Client
    ...
        TenantApi: tenantClient,
```

## Fix

Update to use `TenantAPI` throughout for field and variable names:

```go
type Client struct {
    Api      *api.Client
    TenantAPI tenant.Client
}
...
func NewEnvironmentGroupRuleSetClient(apiClient *api.Client, tenantClient tenant.Client) Client {
    return Client{
        Api:      apiClient,
        TenantAPI: tenantClient,
    }
}
```

---

This issue will be saved in:
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_environment_group_rule_set_naming_low_initialism.md


---

## ISSUE 5

# Title
Naming: Struct and Function Parameter Names Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go

## Problem

The `client` struct and its receiver variable are both named `client`. By convention, struct names in Go should be capitalized (`Client`) when exported, and receiver variable names should preferably be small (e.g., `c`). The current naming can be confusing and does not follow Go idioms.

## Impact

Low/Medium. This impacts code readability and maintainability, especially for people reading or maintaining the code in the future.

## Location

```go
type client struct {
    Api *api.Client
}

func (client *client) ...
```

## Code Issue

```go
type client struct {
    Api *api.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
    //...
}
```

## Fix

Rename the struct to `Client` (capitalize if you want to export), and use a shorter receiver name, typically `c`.

```go
type Client struct {
    Api *api.Client
}

func (c *Client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
    //...
}
```


---

## ISSUE 6

# Naming: Struct Field Name Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The field `Api` in the struct should be named `API` as per Go naming conventions for acronyms.

## Impact

Non-standard naming can reduce code readability and maintainability, especially in large codebases. Severity: **low**.

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

Rename the field to use the all-caps acronym:

```go
type client struct {
	API *api.Client
}
```


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
