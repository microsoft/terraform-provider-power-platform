# Disabled Acceptance Test Without Adequate Justification

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

The `TestAccAnalyticsDataExportsDataSource_Validate_Read` function is an acceptance test that is being skipped due to lack of SP (Service Principal) support. Skipping tests should only be a last resort and should be accompanied by a more explicit reason, logged as an issue or TODO, or better yet, conditionally executed based on available configuration so automated acceptance tests don't silently degrade.

## Impact

Medium. Skipped tests can lead to untested features making their way to production, reducing the reliability of the codebase. Skipped tests without proper visibility or acknowledgment may go unnoticed for long periods, masking gaps in coverage.

## Location

```go
func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	t.Skip("Skipping test due lack of SP support")
	...
}
```

## Code Issue

```go
func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	t.Skip("Skipping test due lack of SP support")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
        ...
```

## Fix

Add a `TODO` comment with a JIRA, GitHub issue, or other visible tracking mechanism, and/or make the skip conditional based on a more robust check, so that the test executes when the underlying capability exists.

```go
func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	// TODO: https://github.com/yourrepo/terraform-provider-power-platform/issues/XXX - Remove skip when SP support is available.
	t.Skip("Skipping test due to lack of SP support. Track at: https://github.com/yourrepo/terraform-provider-power-platform/issues/XXX")

	// Alternatively, make the skip conditional:
	// if !helpers.SupportsSP() {
	//     t.Skip("Skipping test: Service Principal support not available in this environment")
	// }

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
        ...
```
