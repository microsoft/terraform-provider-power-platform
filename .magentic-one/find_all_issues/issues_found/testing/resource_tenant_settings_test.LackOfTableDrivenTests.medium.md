# Title

Lack of Table-driven Testing Patterns

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

Currently, each test is written out with fully repeated steps, configs, and assertions, despite the fact that most tests only vary by a small number of values (e.g., toggling between "true"/"false", or changing thresholds). Go best practices recommend the use of table-driven testing to reduce duplication and make tests more robust, concise, and easier to extend with more scenarios.

## Impact

This reduces test scalability and increases maintenance overhead. To add or tweak test scenarios, contributors must duplicate large code blocks rather than add entries to a concise table. This increases the risk of code rot, logical drift, and authoring errors. Severity: medium.

## Location

Each acceptance/unit test with multiple test steps that only change a few values, e.g., update versus create flows.

## Code Issue

```go
resource.Test(t, resource.TestCase{
    Steps: []resource.TestStep{
        {
            Config: "...",
            Check: ...,
        },
        {
            Config: "...",
            Check: ...,
        },
    },
})
```
(Every step is hand-written and duplicated.)

## Fix

Adopt a Go table-driven style. For instance:

```go
testCases := []struct{
    name   string
    config string
    checks resource.TestCheckFunc
}{
    { name: "create", config: createConfig, checks: createChecks },
    { name: "update", config: updateConfig, checks: updateChecks },
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        resource.Test(t, resource.TestCase{
            Steps: []resource.TestStep{
                {
                    Config: tc.config,
                    Check:  tc.checks,
                },
            },
        })
    })
}
```

Or, if using Steps, at least drive the variations/steps from slices/structs and loop over them.

This will dramatically reduce code bloat, focus the tests, and make it easier to add new test variations.
