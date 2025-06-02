# Title

Lack of Table-driven Testing for Inputs and Variations

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment_test.go

## Problem

All test steps are written as distinct manual blocks. Given the repetition and the number of input permutations, this file would benefit significantly from Go's idiomatic table-driven tests (even within acceptance tests, for structuring `Steps`) to methodically cover variations or edge cases.

## Impact

Low (in this context; medium for maintainability as test matrix grows). While tests work, it increases code length, makes bulk change harder, and may result in copy-paste mistakes.

## Location

Near many multi-step test functions (one manual step per scenario):

```go
resource.Test(t, resource.TestCase{
    Steps: []resource.TestStep{
        { Config: "...", Check: ... },
        { Config: "...", Check: ... },
        // ...
    },
})
```

## Code Issue

```go
Steps: []resource.TestStep{
    { Config: "...", Check: ... },
    { Config: "...", Check: ... },
    // ...
}
```

## Fix

Refactor to something like this:

```go
type testStep struct {
    name   string
    config string
    check  resource.TestCheckFunc
}

var steps = []testStep{
    {name: "disable insights", config: configWithX(false), check: ...},
    {name: "enable insights", config: configWithX(true), check: ...},
    // ...
}

for _, step := range steps {
    t.Run(step.name, func(t *testing.T) {
        resource.Test(t, resource.TestCase{
            Steps: []resource.TestStep{
                {Config: step.config, Check: step.check},
            },
        })
    })
}
```
