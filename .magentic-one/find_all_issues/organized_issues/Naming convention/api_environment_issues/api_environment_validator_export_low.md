# Issue: Anonymous (unexported) validator functions should be exported for consistency

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

There are functions in the file such as `currencyCodeValidator` and `languageCodeValidator` which are not exported (do not start with uppercase) yet perform logic closely related to the exported client methods. In a codebase like this, validator helpers of this kind are typically part of a public contract, and their naming should be consistent for clear documentation and maintainability.

## Impact

- Severity: Low
- Reduces clarity of code purpose and consistency.
- Can create confusion for contributors as to what should be accessible/usable elsewhere in the package, and what is internal-only.

## Location

Function definitions:

```go
func currencyCodeValidator(ctx context.Context, client *api.Client, location string, currencyCode string) error
func languageCodeValidator(ctx context.Context, client *api.Client, location string, languageCode string) error
```

## Code Issue

```go
func currencyCodeValidator(ctx context.Context, client *api.Client, location string, currencyCode string) error
func languageCodeValidator(ctx context.Context, client *api.Client, location string, languageCode string) error
```

## Fix

Make these validator helpers exported to match typical Go conventions for such helpers (unless there is a clear reason for them to be private, and they are only used locally):

```go
func CurrencyCodeValidator(ctx context.Context, client *api.Client, location string, currencyCode string) error
func LanguageCodeValidator(ctx context.Context, client *api.Client, location string, languageCode string) error
```

If they are truly intended to be private, you may add comments explicitly marking them as such.

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_environment_validator_export_low.md`
