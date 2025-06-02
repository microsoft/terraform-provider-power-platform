# Title

Imprecise Aggregate Test Check Structure

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

In both tests, all resource attributes are checked together in a single large `resource.ComposeAggregateTestCheckFunc(...)` block. If any check fails, it's not always clear which specific attribute is problematic, making debugging slower.

## Impact

Low. Readability and debuggability in test output is impacted. The common testing practice is still sound, but maintenance and diagnosis may be a bit harder.

## Location

Lines 28–40 and 71–83.

## Code Issue

```go
Check: resource.ComposeAggregateTestCheckFunc(
	resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", "00000000-0000-0000-0000-000000000001"),
	resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "license_model_type", "StorageDriven"),
	// ... more
),
```

## Fix

Consider splitting checks or grouping related checks into smaller blocks, or at least add comments separating domains being checked (metadata, license details, etc.), and possibly leverage `t.Run` subtests for major sections.

```go
Check: resource.ComposeAggregateTestCheckFunc(
	// Identity attributes
	resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", mockTenantID),

	// License and usage attributes
	resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "license_model_type", "StorageDriven"),
	// ...
),
```

Alternatively, encapsulate repeated patterns in test helper functions.
