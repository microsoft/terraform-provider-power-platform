# Title

Suspicious/Incorrect Domain Constants for RX Cloud

##

/workspaces/terraform-provider-power-platform/internal/constants/constants.go

## Problem

The block defining constants for `RX_*` (presumably related to a national/regional cloud) contains what appear to be copy-paste errors regarding the advisor API domain and scope. The following entries:

```go
RX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.eaglex.ic.gov"
RX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.eaglex.ic.gov/.default"
```

use "eaglex.ic.gov"—the domain used for the EX ("Eagle X") region—rather than an "microsoft.scloud"-based value that would be correct for RX.

## Impact

Medium severity. Using an incorrect or mismatched domain for Power Apps Advisor API for RX-backed tenants results in misconfiguration, which can prevent API calls from working as intended and could cause outages or security holes if the wrong endpoints or tokens are used.

## Location

In the block:

```go
const (
	...
	RX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.eaglex.ic.gov"
	RX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.eaglex.ic.gov/.default"
	...
)
```

## Code Issue

```go
const (
	...
	RX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.eaglex.ic.gov"
	RX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.eaglex.ic.gov/.default"
	...
)
```

## Fix

Change these values to match `microsoft.scloud` instead of using the EX region domain. Correct the constants as follows:

```go
const (
	...
	RX_POWERAPPS_ADVISOR_API_DOMAIN = "api.advisor.powerapps.microsoft.scloud"
	RX_POWERAPPS_ADVISOR_API_SCOPE  = "https://advisor.powerapps.microsoft.scloud/.default"
	...
)
```

This will ensure that requests for the RX Power Apps Advisor API are properly routed and securely authenticated.

---

**Save location:**  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/constants.go-incorrect_domain-medium.md
