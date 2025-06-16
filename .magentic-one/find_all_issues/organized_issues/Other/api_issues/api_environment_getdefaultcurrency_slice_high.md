# Issue: Potential slice out-of-bounds panic when accessing organization settings

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

In the `GetDefaultCurrencyForEnvironment` function, code accesses `orgSettings.Value[0].BaseCurrencyId` directly without checking whether the slice is empty. If `orgSettings.Value` is an empty slice, this will cause a runtime panic.

## Impact

- Severity: High
- Can cause runtime panic and crash the provider.
- Affects reliability and may result in non-obvious bugs in edge/corner cases (e.g., a newly created environment with no settings).

## Location

In `GetDefaultCurrencyForEnvironment`:

```go
values := url.Values{}
values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)
```

## Code Issue

```go
values := url.Values{}
values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)
```

## Fix

Check length of `orgSettings.Value` before accessing its first element:

```go
if len(orgSettings.Value) == 0 {
    return nil, fmt.Errorf("no organization settings found for environment %s", environmentId)
}
values := url.Values{}
values.Add("$filter", "transactioncurrencyid eq "+orgSettings.Value[0].BaseCurrencyId)
```

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_environment_getdefaultcurrency_slice_high.md`
