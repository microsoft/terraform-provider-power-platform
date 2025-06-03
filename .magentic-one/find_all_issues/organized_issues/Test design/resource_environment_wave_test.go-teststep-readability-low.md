# Overuse of anonymous resource.TestStep structs reduces testcase readability

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go

## Problem

The tests invoke `resource.Test` using inline slices of `resource.TestStep` with multi-line raw string literal configs directly inside the slice. While this is functionally correct, it degrades readability and makes the testcases harder to follow, especially as the number of test cases increases. It also makes it difficult to add shared setup or common checks for resources.

## Impact

Low. The code is still valid and runs correctly, but readability suffers, and future maintainability is impaired, especially if you need to share setup or checks between test steps or add new scenarios.

## Location

```go
resource.Test(t, resource.TestCase{
	IsUnitTest:               true,
	ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
	Steps: []resource.TestStep{
		{
			Config: `
			resource "powerplatform_environment_wave" "test" {
				environment_id = "00000000-0000-0000-0000-000000000001"
				feature_name  = "October2024Update"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
				resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
				resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "enabled"),
			),
		},
	},
})
```

## Fix

Assign configs and checkFuncs to well-named variables outside the TestStep struct literals for improved readability:

```go
const createConfig = `
resource "powerplatform_environment_wave" "test" {
	environment_id = "00000000-0000-0000-0000-000000000001"
	feature_name  = "October2024Update"
}`

var createChecks = resource.ComposeTestCheckFunc(
	resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "environment_id", "00000000-0000-0000-0000-000000000001"),
	resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "feature_name", "October2024Update"),
	resource.TestCheckResourceAttr("powerplatform_environment_wave.test", "state", "enabled"),
)

resource.Test(t, resource.TestCase{
	IsUnitTest:               true,
	ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
	Steps: []resource.TestStep{
		{
			Config: createConfig,
			Check:  createChecks,
		},
	},
})
```

This small refactor aids future test extension and improves clarity.
