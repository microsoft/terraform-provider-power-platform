# Title

Hardcoded Tenant ID in Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The test configuration in the `TestUnitTenantCapacityDataSource_Validate_Read` function uses a hardcoded tenant ID (`00000000-0000-0000-0000-000000000001`) multiple times in testing logic. Hardcoding values in test cases reduces flexibility and maintainability, as any change to the tenant ID would require modifications across all references within the code.

## Impact

Hardcoding values can make it difficult to reuse test cases across different environments or tenants. It also increases the likelihood of errors during changes and contributes to technical debt within the test suite. The severity of this issue is **medium**, as it impacts test code maintainability.

## Location

Within the `Config` block of `TestUnitTenantCapacityDataSource_Validate_Read`:

```go
data "powerplatform_tenant_capacity" "capacity" {
	tenant_id = "00000000-0000-0000-0000-000000000001"
}
```

Additionally referenced in assertions:

```go
resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", "00000000-0000-0000-0000-000000000001")
```

## Code Issue

The tenant ID is used directly within the test configuration and assertion, without abstraction or parameterization:

```go
tenant_id = "00000000-0000-0000-0000-000000000001"

resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", "00000000-0000-0000-0000-000000000001")
```

## Fix

Abstract hardcoded values by defining them as constants or environmental variables that can be set externally, thereby improving flexibility and usability.

```go
const DefaultTenantID = "00000000-0000-0000-0000-000000000001"

Config: fmt.Sprintf(`
data "powerplatform_tenant_capacity" "capacity" {
	tenant_id = "%s"
}`, DefaultTenantID),

Check: resource.ComposeAggregateTestCheckFunc(
	resource.TestCheckResourceAttr("data.powerplatform_tenant_capacity.capacity", "tenant_id", DefaultTenantID),
)
```

- Using constants ensures tenant ID modifications are centralized, reducing duplication and promoting cleaner code.
- Tests become easily configurable across different environments, improving maintainability.
