# Potential for Test Flakiness Due to Use of rand.IntN Without Fixed Seed

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

The use of `rand.IntN(9999)` (from `math/rand/v2`) in resource group naming introduces randomness into test resource names. This can result in non-reproducible test runs, which is undesirable for CI/CD automation.

## Impact

- **Reproducibility**: Hard to debug failures.
- **Potential Flakiness**: Unpredictable test IDs.

**Severity: Medium**

## Location

In every test using:

```go
rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
```

## Code Issue

```go
rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
```

## Fix

Seed the random number generator with a constant value (or preferably, with `t.Name()` for per-test reproducibility), or use UUIDs or deterministic suffixes for test resource names.

```go
r := rand.New(rand.NewSource(int64(hash(t.Name()))) )
rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(r.Intn(9999))
```

Or, if randomness isn’t required, use a deterministic fixture naming scheme.
