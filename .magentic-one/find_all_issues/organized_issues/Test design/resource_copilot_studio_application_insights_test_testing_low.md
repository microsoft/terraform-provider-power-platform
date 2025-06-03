# Title

Use of Randomness with rand.IntN Without Fixed Seed Reduces Test Reproducibility

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights_test.go

## Problem

The acceptance/integration test uses `rand.IntN(9999)` without explicitly setting a seed. This makes the test non-deterministic: generated resource names and IDs will be random on every run, making it harder to reproduce failures based on logs or to identify test artifacts left behind in acceptance/integration environments.

## Impact

Low severity for a typical provider, but potentially medium for CI/CD systems and test troubleshooting. Non-deterministic artifacts create challenges for cleaning up resources and tracing issues in persistent test environments.

## Location

Acceptance test function, resource name construction:

```go
name     = "power-platform-app-insights-rg-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
...
name = "power-platform-app-insights-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
```

## Code Issue

```go
name     = "power-platform-app-insights-rg-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
...
name = "power-platform-app-insights-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
```

## Fix

Set the random seed to a fixed or time-dependent value at the start of the test (or via a test utility), or use a deterministic suffix based on test case/ID if reproducibility is important:

```go
// At top of test/func or in package init:
seed := time.Now().UnixNano()
rand.Seed(seed)
t.Logf("Test random seed: %d", seed)
```

Alternatively, to get reproducibility:

```go
const testSeed = 42
rand.Seed(testSeed)
```

Or, derive from the test context:

```go
rand.Seed(int64(hashOfTestName(t.Name())))
```

This makes resource names predictable for a given test run or seed provided to CI.

---

Continuing to check for further issues.