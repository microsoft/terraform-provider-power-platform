# Issue #2: Missing Validation or Tests for Edge Cases and Error Conditions

##

`/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages_test.go`

## Problem

The test case does not include validations for edge cases and scenarios where the input data might fail or return unexpected results. Specifically:
- No test cases for invalid `location` values.
- No test cases to handle null or empty attributes from the data source.
- Absence of mock data that could simulate failures.

## Impact

Without testing edge cases:
- Possible bugs or failures in the production implementation may go unnoticed.
- The reliability and robustness of the data source checks are not fully guaranteed.
- Increased risk of unhandled exceptions in real-world scenarios.

**Severity**: High

## Location

Function: `TestAccLanguagesDataSource_Validate_Read`

## Code Issue

```go
func TestAccLanguagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_languages" "all_languages_for_unitedstates" {
					location = "unitedstates"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.display_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.locale_id", regexp.MustCompile(helpers.StringRegex)),
				),
			},
		},
	})
}
```

## Fix

Introduce additional test steps to handle edge cases and error conditions. Validate scenarios with invalid inputs or responses from the data source.

```go
func TestAccLanguagesDataSource_Validate_EdgeCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_languages" "invalid_languages" {
					location = "invalid-location" // Invalid location simulating incorrect input.
				}`,
				ExpectError: regexp.MustCompile("expected error"), // Check for expected error on invalid location.
			},
			{
				Config: `
				data "powerplatform_languages" "null_location" {
					location = "" // Null or empty location simulation.
				}`,
				ExpectError: regexp.MustCompile("expected error"), // Check for expected error on null location.
			},
			{
				Config: `
				data "powerplatform_languages" "out_of_boundary_languages" {
					location = "someplace" // Simulating unsupported location.
				}`,
				ExpectError: regexp.MustCompile("unsupported location"), // Specific unsupported locations handled correctly.
			},
		},
	})
}
```