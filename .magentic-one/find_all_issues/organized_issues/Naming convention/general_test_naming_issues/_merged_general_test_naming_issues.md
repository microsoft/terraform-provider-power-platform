# General Test Naming Issues - Merged Issues

## ISSUE 1

# Table-driven Test Naming Inconsistent

##

/workspaces/terraform-provider-power-platform/internal/api/client_test.go

## Problem

In the `TestUnitIsCaeChallengeResponse`, table-driven tests are used but the `name` field values are not always precise or following a clear pattern (e.g., missing error details, slight redundancy). While this is not a functional bug, inconsistent table-test names may cause confusion in interpreting failures or understanding intent, especially when tests grow or become more complex. Consistent and descriptive naming in table-driven tests helps readability and maintainability.

## Impact

Low severity; mostly affects readability and the usefulness of test output on failure.

## Location

```go
{
	name: "401 status with WWW-Authenticate header but missing insufficient_claims",
	...
},
// ...
{
	name: "Valid CAE challenge response with complex header",
	...
},
```

## Code Issue

Inconsistent test case naming style.

## Fix

Ensure test `name` fields in table-driven tests are consistently descriptive and use a common style (e.g., start with a status/condition, then mention what is being tested, optionally include expected outcome). For example:

```go
{
	name: "401 Unauthorized with WWW-Authenticate header: missing insufficient_claims",
	...
},
{
	name: "401 Unauthorized with WWW-Authenticate header: valid CAE challenge, complex header",
	...
},
```


---

## ISSUE 2

# Inconsistent Naming in Test Table Types

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The test table struct type for each test is declared as `testData` in all, which is generic and could be more descriptive. Also, naming the struct as `testCase` in the loop (shadowing the table variable) can be confusing.

## Impact

Severity: **Low**  
Minor readability/maintainability issue.

## Location

```go
type testData struct {
    // ...
}
for _, testCase := range []testData{ ... }
```

## Code Issue

```go
type testData struct {
    // ...
}
for _, testCase := range []testData{
    // ...
}
```

## Fix

Rename the struct to a more specific name, e.g. `configStringTestCase`, and use `tc` in the loop:

```go
type configStringTestCase struct {
    // ...
}
for _, tc := range []configStringTestCase{
    // ...
}
```

Repeat for other test tables.


---

## ISSUE 3

# Title

Minor Typo: Function Call `NewBillingPoliciesEnvironmetsDataSource` (Should be "Environments")

##

internal/provider/provider_test.go

## Problem

The function `NewBillingPoliciesEnvironmetsDataSource` contains a probable typo in "Environmets". This should likely be spelled "Environments" for consistency with other naming.

## Impact

Low. This is only a spelling/typo issue, but consistent naming is important for code understanding.

## Location

```go
licensing.NewBillingPoliciesEnvironmetsDataSource(),
```

## Fix

Update the code to use the correct spelling:

```go
licensing.NewBillingPoliciesEnvironmentsDataSource(),
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
