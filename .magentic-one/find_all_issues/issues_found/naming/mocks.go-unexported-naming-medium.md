# Issue: Unexported function and variable names don't comply with Go naming conventions

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

Functions and variables that are not intended to be exported (used externally) should start with a lowercase letter according to Go conventions. Example: `TestUnitTestProtoV6ProviderFactories` and `TestAccProtoV6ProviderFactories` are not exported (used only in tests and possibly mocks), but start with uppercase letters, which is reserved for exported identifiers. Similarly, function names like `TestsEntraLicesingGroupName` are not following the convention.

## Impact

**Severity: Medium**  
Having non-exported elements with uppercase names confuses maintainability and reduces clarity in code navigation. It also makes linters and IDEs treat these objects as potentially exported, which can make refactoring harder and make maintainers unsure about use scope.

## Location

Lines declaring variables and functions such as:

```go
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
...
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
...
}

func TestsEntraLicesingGroupName() string {
...
}
```

## Code Issue

```go
var TestUnitTestProtoV6ProviderFactories = ...
var TestAccProtoV6ProviderFactories = ...
func TestsEntraLicesingGroupName() string { ... }
```

## Fix

Rename unexported identifiers to start with lowercase letters. For example:

```go
var unitTestProtoV6ProviderFactories = ...
var accProtoV6ProviderFactories = ...
func testsEntraLicensingGroupName() string { ... }
```

If these elements are intended to be exported, keep names as-is; otherwise, they should be private unless there’s a clear need for external use.

---

**Save this file as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/mocks.go-unexported-naming-medium.md`
