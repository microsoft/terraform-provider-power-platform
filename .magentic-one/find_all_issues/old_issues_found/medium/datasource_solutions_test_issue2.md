# Title
Missing Unit Test Assertions in `TestUnitSolutionsDataSource_Validate_No_Dataverse`.

## File Path
/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions_test.go

## Problem
The `TestUnitSolutionsDataSource_Validate_No_Dataverse` test case is missing sufficient assertions and validations in the `resource.Test` steps. The `Check` function is defined but does not validate any attributes or solution properties to confirm correctness.

## Impact
This reduces the reliability of the test and diminishes its ability to detect regressions or incorrect behavior. Non-assertive tests provide limited value and can lead to undetected errors in the code.

**Severity: Medium**

## Code Location
```go
resource.Test(t, resource.TestCase{
	IsUnitTest:               true,
	ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
	Steps: []resource.TestStep{
		{
			Config: `
			resource "powerplatform_environment" "env" {
				display_name                              = "displayname"
				location                                  = "europe"
				environment_type                          = "Sandbox"
			}

			data "powerplatform_solutions" "all" {
				environment_id = powerplatform_environment.env.id
			}`,
			ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
			Check:       resource.ComposeAggregateTestCheckFunc(),
		},
	},
})
```

## Fix
Add meaningful assertions in the `Check` function using `TestCheckResourceAttr` or similar constructs to validate that expected errors, behaviors, and API responses are properly tested for correctness.

### Fixed Code
```go
resource.Test(t, resource.TestCase{
	IsUnitTest:               true,
	ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
	Steps: []resource.TestStep{
		{
			Config: `
			resource "powerplatform_environment" "env" {
				display_name                              = "displayname"
				location                                  = "europe"
				environment_type                          = "Sandbox"
			}

			data "powerplatform_solutions" "all" {
				environment_id = powerplatform_environment.env.id
			}`,
			ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.#", "0"),
				resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "id", powerplatform_environment.env.id),
			),
		},
	},
})
```

### Explanation
The fixed code introduces assertions within the `Check` function to validate that the `solutions` dataset is empty and that the `id` attribute of the `data` block matches the environment's ID. This ensures the test properly validates behavior when there is no Dataverse environment.
