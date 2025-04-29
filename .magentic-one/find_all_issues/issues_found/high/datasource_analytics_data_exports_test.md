# Title

Unnecessary Skip of Acceptance Test

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

The `TestAccAnalyticsDataExportsDataSource_Validate_Read` test uses `t.Skip()` to skip acceptance testing due to a lack of service principal support. While such skipping can be appropriate in some situations, it results in unverified functionality for important features. Tests are meant to ensure the integrity of the codebase and skipping them inhibits this goal.

## Impact

Skipping acceptance tests prevents regular validation of critical code behavior and may inadvertently allow unnoticed regressions. Severity: **high**

## Location

File: `datasource_analytics_data_exports_test.go` 

Function: `TestAccAnalyticsDataExportsDataSource_Validate_Read`

## Code Issue

```go
func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	t.Skip("Skipping test due lack of SP support")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: "data \"powerplatform_analytics_data_exports\" \"test\" {}",
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex)),
					extra validation checks omitted...
				),
			},
		},
	})
}
```

## Fix

Enable the test and use proper environmental checks to ensure it runs only in supported environments instead of skipping it outright:

```go
func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	if os.Getenv("SP_SUPPORT") != "true" {
		 t.Skip("Skipping test due to lack of SP support")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: "data \"powerplatform_analytics_data_exports\" \"test\" {}",
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex)),
					extra validation checks omitted...
				),
			},
		},
	})
}
```