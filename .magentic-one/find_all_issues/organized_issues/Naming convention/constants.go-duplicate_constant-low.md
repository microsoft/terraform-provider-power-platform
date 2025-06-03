# Title

Duplicated Constant: `EX_AUTHORITY_HOST` and `EX_OAUTH_AUTHORITY_URL` (and same in RX)

##
/workspaces/terraform-provider-power-platform/internal/constants/constants.go

## Problem

Within the `EX_*` and `RX_*` constant blocks, both `*_AUTHORITY_HOST` and `*_OAUTH_AUTHORITY_URL` are defined and assigned the same value:

In EX:
```go
EX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.eaglex.ic.gov/"
EX_AUTHORITY_HOST      = "https://login.microsoftonline.eaglex.ic.gov/"
```

In RX:
```go
RX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.microsoft.scloud/"
RX_AUTHORITY_HOST      = "https://login.microsoftonline.microsoft.scloud/"
```

This introduces unnecessary duplication. Additionally, the naming difference (`HOST` vs. `URL`) is not meaningful, as both values are URLs (not just hosts).

## Impact

Low severity for most cases, but this can cause confusion, maintenance burden, and the risk of the two diverging in the future. If a fix is made to one, it may be missed in the other.

## Location

```go
const (
	EX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.eaglex.ic.gov/"
	EX_AUTHORITY_HOST      = "https://login.microsoftonline.eaglex.ic.gov/"
)
const (
	RX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.microsoft.scloud/"
	RX_AUTHORITY_HOST      = "https://login.microsoftonline.microsoft.scloud/"
)
```

## Code Issue

```go
EX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.eaglex.ic.gov/"
EX_AUTHORITY_HOST      = "https://login.microsoftonline.eaglex.ic.gov/"
...
RX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.microsoft.scloud/"
RX_AUTHORITY_HOST      = "https://login.microsoftonline.microsoft.scloud/"
```

## Fix

Remove the unnecessary duplicate, and ensure that only one well-named constant exists. Use one of the following approaches:

```go
const (
	EX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.eaglex.ic.gov/"
	// Remove EX_AUTHORITY_HOST
)
const (
	RX_OAUTH_AUTHORITY_URL = "https://login.microsoftonline.microsoft.scloud/"
	// Remove RX_AUTHORITY_HOST
)
```
Or, if a more generic name is required, use a common one for all types (prefer "URL" as the suffix since it's a full URL).

---

**Save location:**  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/constants.go-duplicate_constant-low.md
